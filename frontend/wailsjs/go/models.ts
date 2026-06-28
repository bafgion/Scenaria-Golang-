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
	export class ExportRequest {
	    inputPath: string;
	    output: string;
	    format: string;
	    baseURL: string;
	
	    static createFrom(source: any = {}) {
	        return new ExportRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.inputPath = source["inputPath"];
	        this.output = source["output"];
	        this.format = source["format"];
	        this.baseURL = source["baseURL"];
	    }
	}
	export class ImportRequest {
	    jsonPath: string;
	    outputPath: string;
	
	    static createFrom(source: any = {}) {
	        return new ImportRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.jsonPath = source["jsonPath"];
	        this.outputPath = source["outputPath"];
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
	export class RecordRequest {
	    url: string;
	    output: string;
	    idleSeconds: number;
	    headless: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RecordRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.output = source["output"];
	        this.idleSeconds = source["idleSeconds"];
	        this.headless = source["headless"];
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
	    allureDir: string;
	
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
	        this.allureDir = source["allureDir"];
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
	export class ValidationIssue {
	    line: number;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ValidationIssue(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.line = source["line"];
	        this.message = source["message"];
	    }
	}

}

