from rpi_rf import RFDevice
import signal
from threading import Thread
import time
from queue import Queue
from switch import Switch


class Switch433Controller(Thread):
    def __init__(self, rx_pin, tx_pin, switches):
        super(Switch433Controller, self).__init__()
        self.setDaemon(True)
        self._evt_q = Queue()
        self._switches = switches
        self._rxdevice = RFDevice(rx_pin)
        self._rxdevice.enable_rx()
        self._txdevice = RFDevice(tx_pin)
        self._txdevice.enable_tx()

        for sw in switches:
            switches[sw].set_ctl(self)

        signal.signal(signal.SIGINT, self.exithandler())

    def exithandler(self):
        def exithandler(signal, frame):
            self._rxdevice.cleanup()
            self._txdevice.cleanup()
            sys.exit(0)
        return exithandler

    def send(self, codes, repeat=0):
        if not isinstance(codes, list):
            codes = [codes]
        for code in codes:
            for _ in range(repeat+1):
                self._txdevice.tx_code(code, 1, 293)

    @property
    def evt_queue(self):
        return self._evt_q

    def _check_code(self, code):
        #  print("received code:'{}'".format(code))
        for sw in self._switches:
            if code == self._switches[sw].ON:
                print("Switch " + sw + " turned on")
                self.evt_queue.put((sw, True))
            elif code == self._switches[sw].OFF:
                print("Switch " + sw + " turned off")
                self.evt_queue.put((sw, False))



    def run(self):
        timestamp = None
        print("Listening for codes")
        while True:
            if self._rxdevice.rx_code_timestamp != timestamp:
                timestamp = self._rxdevice.rx_code_timestamp
                self._check_code(self._rxdevice.rx_code)
            time.sleep(1)


class Switch433(Switch):
    def __init__(self, name, state, oncode, offcode):
        super(Switch433, self).__init__(name, state)
        self._on = oncode
        self._off = offcode

    def set_ctl(self, ctl):
        self._ctl = ctl

    @property
    def ON(self):
        return self._on

    @property
    def OFF(self):
        return self._off

    def turn_on(self):
        self._ctl.send(self._on)

    def turn_off(self):
        self._ctl.send(self._off)
