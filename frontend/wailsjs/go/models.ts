export namespace consts {
	
	export enum EVENT {
	    EVENT_PROGRESS = "EVENT_PROGRESS",
	    EVENT_RESULT = "EVENT_RESULT",
	}

}

export namespace model {
	
	export class CrackResult {
	    ID: number;
	    TaskID: string;
	    Target: string;
	    Service: string;
	    Username: string;
	    Password: string;
	
	    static createFrom(source: any = {}) {
	        return new CrackResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.TaskID = source["TaskID"];
	        this.Target = source["Target"];
	        this.Service = source["Service"];
	        this.Username = source["Username"];
	        this.Password = source["Password"];
	    }
	}
	export class CrackTask {
	    Status: number;
	    StartTime: number;
	    EndTime: number;
	    Index: number;
	    Total: number;
	    Progress: number;
	    LastCost: number;
	    ID: string;
	    Targets: string;
	    Usernames: string;
	    Passwords: string;
	    Proxies: string;
	    Thread: number;
	    Interval: number;
	    MaxRuntime: number;
	
	    static createFrom(source: any = {}) {
	        return new CrackTask(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Status = source["Status"];
	        this.StartTime = source["StartTime"];
	        this.EndTime = source["EndTime"];
	        this.Index = source["Index"];
	        this.Total = source["Total"];
	        this.Progress = source["Progress"];
	        this.LastCost = source["LastCost"];
	        this.ID = source["ID"];
	        this.Targets = source["Targets"];
	        this.Usernames = source["Usernames"];
	        this.Passwords = source["Passwords"];
	        this.Proxies = source["Proxies"];
	        this.Thread = source["Thread"];
	        this.Interval = source["Interval"];
	        this.MaxRuntime = source["MaxRuntime"];
	    }
	}

}

