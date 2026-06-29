export namespace gui {
	
	export class AppSettingsDTO {
	    browser: string;
	    headless: boolean;
	    parallelWorkers: number;
	    maxLoopIterations: number;
	    filterRecording: boolean;
	    navOnlyRecording: boolean;
	    hoverRecord: boolean;
	    toolbarCompact: boolean;
	    stepsPanelVisible: boolean;
	    stepsPanelHeight: number;
	    sidebarWidth: number;
	    recentProjects: string[];
	    recentFeatures: string[];
	    checkUpdatesOnStartup: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AppSettingsDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.browser = source["browser"];
	        this.headless = source["headless"];
	        this.parallelWorkers = source["parallelWorkers"];
	        this.maxLoopIterations = source["maxLoopIterations"];
	        this.filterRecording = source["filterRecording"];
	        this.navOnlyRecording = source["navOnlyRecording"];
	        this.hoverRecord = source["hoverRecord"];
	        this.toolbarCompact = source["toolbarCompact"];
	        this.stepsPanelVisible = source["stepsPanelVisible"];
	        this.stepsPanelHeight = source["stepsPanelHeight"];
	        this.sidebarWidth = source["sidebarWidth"];
	        this.recentProjects = source["recentProjects"];
	        this.recentFeatures = source["recentFeatures"];
	        this.checkUpdatesOnStartup = source["checkUpdatesOnStartup"];
	    }
	}
	export class EditorStepRow {
	    line: number;
	    keyword: string;
	    text: string;
	    action: string;
	    element: string;
	    value: string;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new EditorStepRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.line = source["line"];
	        this.keyword = source["keyword"];
	        this.text = source["text"];
	        this.action = source["action"];
	        this.element = source["element"];
	        this.value = source["value"];
	        this.error = source["error"];
	    }
	}
	export class ScenarioHintDTO {
	    id: string;
	    title: string;
	    stepIndex: number;
	    severity: string;
	    autoFixable: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ScenarioHintDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.stepIndex = source["stepIndex"];
	        this.severity = source["severity"];
	        this.autoFixable = source["autoFixable"];
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
	export class ExportPreview {
	    stepCount: number;
	    scenarioTitle: string;
	    issues: ValidationIssue[];
	    hints: ScenarioHintDTO[];
	
	    static createFrom(source: any = {}) {
	        return new ExportPreview(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stepCount = source["stepCount"];
	        this.scenarioTitle = source["scenarioTitle"];
	        this.issues = this.convertValues(source["issues"], ValidationIssue);
	        this.hints = this.convertValues(source["hints"], ScenarioHintDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class ExportRequest {
	    inputPath: string;
	    output: string;
	    format: string;
	    baseURL: string;
	    force: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ExportRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.inputPath = source["inputPath"];
	        this.output = source["output"];
	        this.format = source["format"];
	        this.baseURL = source["baseURL"];
	        this.force = source["force"];
	    }
	}
	export class HTTPAuthCredentials {
	    username: string;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new HTTPAuthCredentials(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.password = source["password"];
	    }
	}
	export class HTTPAuthRequest {
	    host: string;
	    username: string;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new HTTPAuthRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.host = source["host"];
	        this.username = source["username"];
	        this.password = source["password"];
	    }
	}
	export class HighlightSpan {
	    text: string;
	    kind: string;
	
	    static createFrom(source: any = {}) {
	        return new HighlightSpan(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.text = source["text"];
	        this.kind = source["kind"];
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
	export class PickSelectorResult {
	    selector: string;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new PickSelectorResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.selector = source["selector"];
	        this.error = source["error"];
	    }
	}
	export class PickerStepChoice {
	    label: string;
	    stepBody: string;
	    description: string;
	    preview: string;
	
	    static createFrom(source: any = {}) {
	        return new PickerStepChoice(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.label = source["label"];
	        this.stepBody = source["stepBody"];
	        this.description = source["description"];
	        this.preview = source["preview"];
	    }
	}
	export class PluginEntryDTO {
	    name: string;
	    source: string;
	    installedAt: string;
	    id: string;
	    description: string;
	    commands: string[];
	    runnable: boolean;
	    vanessa: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PluginEntryDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.source = source["source"];
	        this.installedAt = source["installedAt"];
	        this.id = source["id"];
	        this.description = source["description"];
	        this.commands = source["commands"];
	        this.runnable = source["runnable"];
	        this.vanessa = source["vanessa"];
	    }
	}
	export class ProjectArtifacts {
	    allureDir: string;
	    tracesDir: string;
	    videosDir: string;
	    htmlReport: string;
	    junitReport: string;
	
	    static createFrom(source: any = {}) {
	        return new ProjectArtifacts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.allureDir = source["allureDir"];
	        this.tracesDir = source["tracesDir"];
	        this.videosDir = source["videosDir"];
	        this.htmlReport = source["htmlReport"];
	        this.junitReport = source["junitReport"];
	    }
	}
	export class PluginRunRequest {
	    name: string;
	    dryRun: boolean;
	    tag: string;
	    excludeTags: string[];
	    scenario: string;
	
	    static createFrom(source: any = {}) {
	        return new PluginRunRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.dryRun = source["dryRun"];
	        this.tag = source["tag"];
	        this.excludeTags = source["excludeTags"];
	        this.scenario = source["scenario"];
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
	export class ProjectReplaceRequest {
	    find: string;
	    replace: string;
	    caseSensitive: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ProjectReplaceRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.find = source["find"];
	        this.replace = source["replace"];
	        this.caseSensitive = source["caseSensitive"];
	    }
	}
	export class ProjectReplaceResult {
	    filesChanged: number;
	    replacements: number;
	    files: string[];
	
	    static createFrom(source: any = {}) {
	        return new ProjectReplaceResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filesChanged = source["filesChanged"];
	        this.replacements = source["replacements"];
	        this.files = source["files"];
	    }
	}
	export class RecentsDTO {
	    projects: string[];
	    features: string[];
	
	    static createFrom(source: any = {}) {
	        return new RecentsDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.projects = source["projects"];
	        this.features = source["features"];
	    }
	}
	export class RecordRequest {
	    url: string;
	    output: string;
	    idleSeconds: number;
	    headless: boolean;
	    filterRecording: boolean;
	    navOnlyRecording: boolean;
	    hoverRecord: boolean;
	    appendTo: string;
	    testClient: string;
	
	    static createFrom(source: any = {}) {
	        return new RecordRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.output = source["output"];
	        this.idleSeconds = source["idleSeconds"];
	        this.headless = source["headless"];
	        this.filterRecording = source["filterRecording"];
	        this.navOnlyRecording = source["navOnlyRecording"];
	        this.hoverRecord = source["hoverRecord"];
	        this.appendTo = source["appendTo"];
	        this.testClient = source["testClient"];
	    }
	}
	export class RefactorResult {
	    text: string;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new RefactorResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.text = source["text"];
	        this.count = source["count"];
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
	    traceDir: string;
	    videoDir: string;
	    htmlPath: string;
	    junitPath: string;
	    targets: string[];
	    browser: string;
	    workers: number;
	    slowMo: number;
	
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
	        this.traceDir = source["traceDir"];
	        this.videoDir = source["videoDir"];
	        this.htmlPath = source["htmlPath"];
	        this.junitPath = source["junitPath"];
	        this.targets = source["targets"];
	        this.browser = source["browser"];
	        this.workers = source["workers"];
	        this.slowMo = source["slowMo"];
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
	export class RunResultEntry {
	    path: string;
	    success: boolean;
	    message: string;
	    runner: string;
	    at: string;
	
	    static createFrom(source: any = {}) {
	        return new RunResultEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.success = source["success"];
	        this.message = source["message"];
	        this.runner = source["runner"];
	        this.at = source["at"];
	    }
	}
	
	export class ScenarioHintFixRequest {
	    text: string;
	    hintId: string;
	    stepIndex: number;
	
	    static createFrom(source: any = {}) {
	        return new ScenarioHintFixRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.text = source["text"];
	        this.hintId = source["hintId"];
	        this.stepIndex = source["stepIndex"];
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

