'use strict';



class ControlPage {

    constructor() {
        this.data = null
        this.tournament = 0
        this.table = null

    }

    refresh(control, id) {
        control.loadData(id);
    }

}
