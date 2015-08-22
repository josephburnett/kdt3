#    Copyright 2009 Google Inc.
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#         http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.

"""Simple web application load testing script.

This is a simple web application load
testing skeleton script. Modify the code between !!!!!
to make the requests you want load tested.
"""


import httplib2
import random
import socket
import time
from threading import Event
from threading import Thread
from threading import current_thread
from urllib import urlencode

# Modify these values to control how the testing is done

# How many threads should be running at peak load.
NUM_THREADS = 1

# How many minutes the test should run with all threads active.
TIME_AT_PEAK_QPS = 1 # minutes

# How many seconds to wait between starting threads.
# Shouldn't be set below 30 seconds.
DELAY_BETWEEN_THREAD_START = 30 # seconds

# Which endpoint should we hit?
ENDPOINT = "http://localhost:8080"

quitevent = Event()

def assertRequest(h, expected):
    """Make a request and assert the expected response."""
    resp, cont = h.request(expected['url'])
    if resp.status != expected['status']:
        print "Unexpected status: " + resp.status
        return False
    if not resp['content-location'].startswith(expected['loc']):
        print "Unexpected content-location: " + resp['content-location']
        return False
    return True

def runGame(h, requestList):
    """Make all listed requests and assert they all succeed."""
    for r in requestList:
        if not assertRequest(h, r): return False
    return True

def threadproc():
    """This function is executed by each thread."""
    print "Thread started: %s" % current_thread().getName()
    succeeded = 0
    attempted = 0
    h = httplib2.Http(timeout=30)
    while not quitevent.is_set():
        try:
            attempted += 1
            success = runGame(h, [
                {'url': ENDPOINT, 'status': 200, 'loc': ENDPOINT},
                {'url': ENDPOINT+"/game?handle=Player%201;playerCount=2;k=2;size=3;inarow=3", 'status': 200, 'loc': ENDPOINT+"/game"}
                # Play the rest of a game ...
            ])
            if success:
                succeeded += 1
        except socket.timeout:
            pass

    print "Thread finished: %s (%s/%s)" % (current_thread().getName(), succeeded, attempted)


if __name__ == "__main__":
    runtime = (TIME_AT_PEAK_QPS * 60 + DELAY_BETWEEN_THREAD_START * NUM_THREADS)
    print "Total runtime will be: %d seconds" % runtime
    threads = []
    try:
        for i in range(NUM_THREADS):
            t = Thread(target=threadproc)
            t.start()
            threads.append(t)
            time.sleep(DELAY_BETWEEN_THREAD_START)
        print "All threads running"
        time.sleep(TIME_AT_PEAK_QPS*60)
        print "Completed full time at peak qps, shutting down threads"
    except:
        print "Exception raised, shutting down threads"

    quitevent.set()
    time.sleep(3)
    for t in threads:
        t.join(1.0)
    print "Finished"
