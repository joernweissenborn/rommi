class Switch(object):
    def __init__(self, name, state):
        self._name = name
        self._state = state

    @property
    def name(self):
        return self._name

    @property
    def state(self):
        return self._state

    @state.setter
    def state(self, value):
        self._state = value

    def turn_on(self):
        return NotImplementedError

    def turn_off(self):
        return NotImplementedError


class Switches(object):
    def __init__(self, switches, output):
        self._switches = switches
        self._output = output
        self.update_states_property()

    def turn_switch(self, name, state):
        if state:
            self._switches[name].turn_on()
        else:
            self._switches[name].turn_off()

    def set_state(self, name, state):
        self._switches[name].state = state

    def update_states_property(self):
        names, states = self.states()
        self._output.set_property("States", {"Name": names, "State": states})

    def states(self):
        names = [name for name in self._switches]
        states = [self._switches[name].state for name in self._switches]
        return names, states

    def set_switch(self, request):
        name, state = request.parameter()['Name'], request.parameter()['State']
        print("setting switch '{}': {}".format(name,state))
        self.turn_switch(name, state)
        request.reply({})
        self.update_states_property()


    def set_switch_state(self, name, state):
        print("setting switch '{}': {}".format(name,state))
        self.set_state(name, state)
        self.update_states_property()

