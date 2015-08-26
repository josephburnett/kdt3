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
import sys
import time
from threading import Event
from threading import Thread
from threading import current_thread
from urllib import urlencode
from lxml import html
from string import Template

# Modify these values to control how the testing is done

# How many threads should be running at peak load.
NUM_THREADS = 5

# How many minutes the test should run with all threads active.
TIME_AT_PEAK_QPS = 10 # minutes

# How many seconds to wait between starting threads.
# Shouldn't be set below 30 seconds.
DELAY_BETWEEN_THREAD_START = 30 # seconds

# Which endpoint should we hit?
ENDPOINT = "http://kdtictactoe-1048.appspot.com"

quitevent = Event()

def runStaticGame(h):
    """Play a predetermined 2D-2P 3x3 game."""
    return runGame(h, [
        {'url': ENDPOINT, 'status': 200, 'loc': ENDPOINT},
        {'url': ENDPOINT+"/game?handle=Player%201;playerCount=2;k=5;size=5;inarow=3", 'status': 200, 'loc': ENDPOINT+"/game",
         'parameterize': lambda resp, cont:
            {
                'gameId': gameIds(cont)[0],
                'player1': gameIds(cont)[1][0],
                'player2': gameIds(cont)[1][1]
            }
        },
        {'url': ENDPOINT+"/game/${gameId}?player=${player1}", 'status': 200, 'loc': ENDPOINT+"/game/${gameId}",
         'assertions': [ lambda _, cont: "(your turn)" in cont ]},
        {'url': ENDPOINT+"/move/${gameId}?player=${player1};point=0,0,0,0,0", 'status': 200, 'loc': ENDPOINT+"/game/${gameId}",
         'assertions': [ lambda _, cont: "Move accepted." in cont ]},
        {'url': ENDPOINT+"/game/${gameId}?player=${player2}", 'status': 200, 'loc': ENDPOINT+"/game/${gameId}"},
        {'url': ENDPOINT+"/move/${gameId}?player=${player2};point=0,0,0,1,1", 'status': 200, 'loc': ENDPOINT+"/game/${gameId}"},
        {'url': ENDPOINT+"/game/${gameId}?player=${player1}", 'status': 200, 'loc': ENDPOINT+"/game/${gameId}"},
        {'url': ENDPOINT+"/move/${gameId}?player=${player1};point=0,0,0,1,0", 'status': 200, 'loc': ENDPOINT+"/game/${gameId}"},
        {'url': ENDPOINT+"/game/${gameId}?player=${player2}", 'status': 200, 'loc': ENDPOINT+"/game/${gameId}"},
        {'url': ENDPOINT+"/move/${gameId}?player=${player2};point=0,0,0,2,0", 'status': 200, 'loc': ENDPOINT+"/game/${gameId}"},
        {'url': ENDPOINT+"/game/${gameId}?player=${player1}", 'status': 200, 'loc': ENDPOINT+"/game/${gameId}"},
        {'url': ENDPOINT+"/move/${gameId}?player=${player1};point=0,0,0,0,2", 'status': 200, 'loc': ENDPOINT+"/game/${gameId}"},
        {'url': ENDPOINT+"/game/${gameId}?player=${player2}", 'status': 200, 'loc': ENDPOINT+"/game/${gameId}"},
        {'url': ENDPOINT+"/move/${gameId}?player=${player2};point=0,0,0,2,1", 'status': 200, 'loc': ENDPOINT+"/game/${gameId}"},
        {'url': ENDPOINT+"/game/${gameId}?player=${player2}", 'status': 200, 'loc': ENDPOINT+"/game/${gameId}",
         'assertions': [ lambda _, cont: not "(your turn)" in cont ]},
        {'url': ENDPOINT+"/move/${gameId}?player=${player2};point=0,0,0,2,2", 'status': 200, 'loc': ENDPOINT+"/game/${gameId}",
         'assertions': [ lambda _, cont: "Invalid move: out of turn." in cont ]},
        {'url': ENDPOINT+"/game/${gameId}?player=${player1}", 'status': 200, 'loc': ENDPOINT+"/game/${gameId}",
         'assertions': [ lambda _, cont: "(your turn)" in cont ]},
        {'url': ENDPOINT+"/move/${gameId}?player=${player1};point=0,0,0,0,1", 'status': 200, 'loc': ENDPOINT+"/game/${gameId}",
         'assertions': [ lambda _, cont: "Game over!" in cont,
                         lambda _, cont: "(your turn)" in cont ]},
        {'url': ENDPOINT+"/game/${gameId}?player=${player2}", 'status': 200, 'loc': ENDPOINT+"/game/${gameId}",
         'assertions': [ lambda _, cont: "Game over!" in cont,
                         lambda _, cont: not "(your turn)" in cont ]},
    ])


def assertRequest(h, r, ctx):
    """Make a request and assert the expected response."""
    url = Template(r['url']).substitute(ctx)
    loc = Template(r['loc']).substitute(ctx)
    resp, cont = h.request(url)
    if resp.status != r['status']:
        print "Unexpected status: " + str(resp.status)
        print "url: " + url
        return False, "", ""
    if not resp['content-location'].startswith(loc):
        print "Unexpected content-location: " + resp['content-location']
        print "url: " + url
        return False, "", ""
    if r.has_key('assertions'):
        for a in r['assertions']:
            if not a(resp, cont):
                print "Assertion failed." # TODO: say which one
                return False, "", ""
    return True, resp, cont

def runGame(h, requestList):
    """Make all listed requests and assert they all succeed."""
    ctx = {}
    for r in requestList:
        success, resp, cont = assertRequest(h, r, ctx)
        if not success: return False
        if r.has_key('parameterize'):
            ctx.update(r['parameterize'](resp, cont))
    return True

def gameIds(cont):
    """Read game and player ids from game creation HTML."""
    root = html.fromstring(cont)
    links = root.xpath('//a/@href')
    if len(links) < 2:
        raise Exception("Expected at least 2 links")
    # /game/YZWGN2VSA733XGVXDAVSTCCMLM?player=SCZVA2IHJNXIJBWSL22MGUOCJQ;...
    #       ^                        ^        ^                        ^
    gameId = links[0][5:31]
    playerIds = []
    for link in links:
        playerIds.append(link[39:65])
    return gameId, playerIds

def threadproc():
    """This function is executed by each thread."""
    print "Thread started: %s" % current_thread().getName()
    succeeded = 0
    attempted = 0
    h = httplib2.Http(timeout=30)
    while not quitevent.is_set():
        try:
            attempted += 1
            if runStaticGame(h):
                succeeded += 1
        except socket.timeout:
            pass

    print "Thread finished: %s (%s/%s)" % (current_thread().getName(), succeeded, attempted)

def loadTest():
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

def releaseTest():
    print "Running release tests"
    h = httplib2.Http(timeout=30)
    start = time.time()
    if runStaticGame(h):
        end = time.time()
        print "Static game passed"
        print "Time: "+ str(end-start)
    else:
        print "Static game failed"
        exit(1)
    print "Finished"

if __name__ == "__main__":
    if len(sys.argv) == 2:
        if sys.argv[1] == "load":
            loadTest()
            exit(0)
        if sys.argv[1] == "release":
            releaseTest()
            exit(0)
    print "Usage: python "+sys.argv[0]+" load|release"
    exit(1)
