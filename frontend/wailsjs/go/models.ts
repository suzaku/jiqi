export namespace k8s {
	
	export class NodeCapacity {
	    cpu: number;
	    memory: number;
	
	    static createFrom(source: any = {}) {
	        return new NodeCapacity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cpu = source["cpu"];
	        this.memory = source["memory"];
	    }
	}
	export class NodeUsage {
	    cpu: number;
	    memory: number;
	
	    static createFrom(source: any = {}) {
	        return new NodeUsage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cpu = source["cpu"];
	        this.memory = source["memory"];
	    }
	}
	export class Node {
	    name: string;
	    labels: {[key: string]: string};
	    consolePageURL: string;
	    dashboardURL: string;
	    instanceType: string;
	    usage: NodeUsage;
	    capacity: NodeCapacity;
	
	    static createFrom(source: any = {}) {
	        return new Node(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.labels = source["labels"];
	        this.consolePageURL = source["consolePageURL"];
	        this.dashboardURL = source["dashboardURL"];
	        this.instanceType = source["instanceType"];
	        this.usage = this.convertValues(source["usage"], NodeUsage);
	        this.capacity = this.convertValues(source["capacity"], NodeCapacity);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

export namespace main {
	
	export class QueriedNodes {
	    nodes: k8s.Node[];
	    labels: {[key: string]: string[]};
	
	    static createFrom(source: any = {}) {
	        return new QueriedNodes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nodes = this.convertValues(source["nodes"], k8s.Node);
	        this.labels = source["labels"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

