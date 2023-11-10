export namespace main {
	
	export class Node {
	    name: string;
	    consolePageURL: string;
	    dashboardURL: string;
	
	    static createFrom(source: any = {}) {
	        return new Node(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.consolePageURL = source["consolePageURL"];
	        this.dashboardURL = source["dashboardURL"];
	    }
	}

}

