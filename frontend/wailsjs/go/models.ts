export namespace main {
	
	export class FileItem {
	    name: string;
	    type: string;
	    path: string;
	    source: string;
	    url: string;
	    imageUrl: string;
	    tags: string[];
	
	    static createFrom(source: any = {}) {
	        return new FileItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.path = source["path"];
	        this.source = source["source"];
	        this.url = source["url"];
	        this.imageUrl = source["imageUrl"];
	        this.tags = source["tags"];
	    }
	}

}

