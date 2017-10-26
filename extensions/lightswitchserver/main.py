from pythingiverseio import Input, Output
from switch import Switches
from switch433 import Switch433Controller, Switch433
import threading
from collections import OrderedDict
import umsgpack
import time

descriptor = '''
function SetSwitch(Name string, State bool)
function SetSwitches(State bool)
property States: Name []string, []State bool
'''

descriptor_extension = '''
function RegisterExtension(Extension bin)
'''

GPIO_TX = 17
GPIO_RX = 27


# SW	ON		OFF
# ===============================
# A	1115473		1115476
# B	1118545		1118548
# C	1119333     1119316
# D	1119505		1119508

extension = umsgpack.packb({
    "Name": "LightSwitchServer",
    "Descriptor": descriptor,
    "Actions": [
        {"Name": "Turn On Lights",
         "Function": "SetSwitches",
         "Parameter": umsgpack.packb({"State":True}),
         "Sentences": [
             "licht an",
             "mache das licht an",
             "lichter an",
             "mache die lichter an",
         ]},
        {"Name": "Turn Off Lights",
         "Function": "SetSwitches",
         "Parameter": umsgpack.packb({"State":False}),
         "Sentences": [
             "licht aus",
             "mache das licht aus",
             "lichter aus",
             "mache die lichter aus",
         ]},
        {"Name": "Turn On Vorderlicht",
         "Function": "SetSwitch",
         "Parameter": umsgpack.packb({"Name": "vorderlicht", "State": True}),
         "Sentences": [
             "vorderlicht an",
             "mache das vorderlicht an",
         ]},
        {"Name": "Turn Off Vorderlicht",
         "Function": "SetSwitch",
         "Parameter": umsgpack.packb({"Name": "vorderlicht", "State": False}),
         "Sentences": [
             "vorderlicht aus",
             "mache das vorderlicht aus",
         ]},
        {"Name": "Turn On Lights Hauptlicht",
         "Function": "SetSwitch",
         "Parameter": umsgpack.packb({"Name": "hauptlicht", "State": True}),
         "Sentences": [
             "hauptlicht an",
             "mache das hauptlicht an",
         ]},
        {"Name": "Turn Off Lights hauptlicht",
         "Function": "SetSwitch",
         "Parameter": umsgpack.packb({"Name": "hauptlicht", "State": False}),
         "Sentences": [
             "hauptlicht aus",
             "mache das hauptlicht aus",
         ]},
        {"Name": "Turn On Lights hinterlicht",
         "Function": "SetSwitch",
         "Parameter": umsgpack.packb({"Name": "hinterlicht", "State": True}),
         "Sentences": [
             "hinterlicht an",
             "mache das hinterlicht an",
         ]},
        {"Name": "Turn Off Lights hinterlicht",
         "Function": "SetSwitch",
         "Parameter": umsgpack.packb({"Name": "hinterlicht", "State": False}),
         "Sentences": [
             "hinterlicht aus",
             "mache das hinterlicht aus",
         ]},
    ]
})

SWITCHES433 = OrderedDict()
SWITCHES433["vorderlicht"]= Switch433("A", False, 1115473, 1115476)
SWITCHES433["hauptlicht"]= Switch433("B", False, 1118545, 1118548)
SWITCHES433["hinterlicht"]= Switch433("D", False, 1119505, 1119508)

class ExtenstionInput(threading.Thread):
    def __init__(self):
        self._input = Input(descriptor_extension)
        super(RequestResponder, self).__init__()

    def run(self):
        connected = self._input.connected()

        if connected:
            self._input.trigger_all("RegisterExtension", {"Extension":extension})

        while True:
            if self._input.connected() and not connected:
                self._input.trigger_all("RegisterExtension", {"Extension":extension})
            connected = self._input.connected()
            time.sleep(1)


class RequestResponder(threading.Thread):
    def __init__(self, output, switches):
        self._output = output
        self._switches = switches
        super(RequestResponder, self).__init__()

    def run(self):
        while True:
            request = self._output.get_request()

            print(request.function())

            if request.function() == "SetSwitch":
                self._switches.set_switch(request)

class SwitchEventResponder(threading.Thread):
    def __init__(self, ctl, switches):
        self._ctl = ctl
        self._switches = switches
        super(SwitchEventResponder, self).__init__()

    def run(self):
        while True:
            (name, state) = self._ctl.evt_queue.get()
            self._switches.set_switch_state(name, state)

def main():
    print("Lightserver starting up")
    ctl = Switch433Controller(GPIO_RX, GPIO_TX, SWITCHES433)
    ctl.start()
    o = Output(descriptor)

    switches = Switches(SWITCHES433, o)

    RequestResponder(o, switches).start()
    SwitchEventResponder(ctl, switches).start()
    ExtenstionInput().start()



if __name__ == "__main__":
    main()
