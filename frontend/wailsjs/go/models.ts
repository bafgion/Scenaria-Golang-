export namespace gui {
	
	export class AppSettingsDTO {
	    browser: string;
	    headless: boolean;
	    parallelWorkers: number;
	    maxLoopIterations: number;
	
	    static createFrom(source: any = {}) {
	        return new AppSettingsDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.browser = source["browser"];
	        this.headless = source["headless"];
	        this.parallelWorkers = source["parallelWorkers"];
	        this.maxLoopIterations = source["maxLoopIterations"];
	    }
	}
	export class ProjectInfo {
	    path: string;
	    features: string[];
	    tags: string[];
	
	    static createFrom(source: any = {}) {
	        return new ProjectInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.features = source["features"];
	        this.tags = source["tags"];
	    }
	}
	export class RunRequest {
	    tag: string;
	    testClient: string;
	    vars: Record<string, string>;
	    dryRun: boolean;
	    headed: boolean;
	    engine: string;
	    installPlaywright: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RunRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tag = source["tag"];
	        this.testClient = source["testClient"];
	        this.vars = source["vars"];
	        this.dryRun = source["dryRun"];
	        this.headed = source["headed"];
	        this.engine = source["engine"];
	        this.installPlaywright = source["installPlaywright"];
	    }
	}
	export class RunResult {
	    output: string;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new RunResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.output = source["output"];
	        this.error = source["error"];
	    }
	}
	export class StepCatalogEntry {
	    category: string;
	    template: string;
	    help: string;
	
	    static createFrom(source: any = {}) {
	        return new StepCatalogEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.category = source["category"];
	        this.template = source["template"];
	        this.help = source["help"];
	    }
	}

}

