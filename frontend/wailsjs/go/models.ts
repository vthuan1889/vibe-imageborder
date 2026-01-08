export namespace models {
	
	export class ProcessRequest {
	    productImages: string[];
	    frameImage: string;
	    templatePath: string;
	    fieldValues: Record<string, string>;
	    outputDir: string;
	    format: string;
	    quality: number;
	
	    static createFrom(source: any = {}) {
	        return new ProcessRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.productImages = source["productImages"];
	        this.frameImage = source["frameImage"];
	        this.templatePath = source["templatePath"];
	        this.fieldValues = source["fieldValues"];
	        this.outputDir = source["outputDir"];
	        this.format = source["format"];
	        this.quality = source["quality"];
	    }
	}

}

export namespace updater {
	
	export class UpdateInfo {
	    available: boolean;
	    current: string;
	    latest: string;
	    downloadUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.available = source["available"];
	        this.current = source["current"];
	        this.latest = source["latest"];
	        this.downloadUrl = source["downloadUrl"];
	    }
	}

}

