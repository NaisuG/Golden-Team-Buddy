export namespace catalog {
	
	export class Champion {
	    key: string;
	    name: string;
	    cost: number;
	    traits: string[];
	
	    static createFrom(source: any = {}) {
	        return new Champion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.name = source["name"];
	        this.cost = source["cost"];
	        this.traits = source["traits"];
	    }
	}

}

export namespace engine {
	
	export class ActiveTrait {
	    Name: string;
	    Count: number;
	    Style: string;
	
	    static createFrom(source: any = {}) {
	        return new ActiveTrait(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Count = source["Count"];
	        this.Style = source["Style"];
	    }
	}
	export class Score {
	    TraitStrength: number;
	    BoardCount: number;
	    Synergy: number;
	
	    static createFrom(source: any = {}) {
	        return new Score(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TraitStrength = source["TraitStrength"];
	        this.BoardCount = source["BoardCount"];
	        this.Synergy = source["Synergy"];
	    }
	}
	export class Variant {
	    Champions: string[];
	    Traits: ActiveTrait[];
	    Score: Score;
	    Children: Variant[];
	
	    static createFrom(source: any = {}) {
	        return new Variant(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Champions = source["Champions"];
	        this.Traits = this.convertValues(source["Traits"], ActiveTrait);
	        this.Score = this.convertValues(source["Score"], Score);
	        this.Children = this.convertValues(source["Children"], Variant);
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

}

