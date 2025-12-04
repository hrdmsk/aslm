export namespace db {
	
	export class ProductInfo {
	    Path: string;
	    Name: string;
	    Url: string;
	    ImageUrl: string;
	    Tags: string[];
	
	    static createFrom(source: any = {}) {
	        return new ProductInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Path = source["Path"];
	        this.Name = source["Name"];
	        this.Url = source["Url"];
	        this.ImageUrl = source["ImageUrl"];
	        this.Tags = source["Tags"];
	    }
	}

}

export namespace main {
	
	export class BoothInfo {
	    productUrl: string;
	    imageUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new BoothInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.productUrl = source["productUrl"];
	        this.imageUrl = source["imageUrl"];
	    }
	}
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

