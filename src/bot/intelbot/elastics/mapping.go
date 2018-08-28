package elastics

const StatusMapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"status":{
			"properties":{
				"id":{
					"type":"keyword"
				},
				"server":{
					"type":"keyword"
				},
				"created_at":{
					"type":"date"
				},
				"sensitive":{
					"type":"boolean"
				},
				"account_id":{
					"type":"keyword"
				},
				"content":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"reblogs_count":{
					"type":"long"
				},
				"favourites_count":{
					"type":"long"
				}
			}
		}
	}
}`

const LocalMapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"local":{
			"properties":{
				"id":{
					"type":"keyword"
				},
				"server":{
					"type":"keyword"
				},
				"created_at":{
					"type":"date"
				},
				"sensitive":{
					"type":"boolean"
				},
				"account_id":{
					"type":"keyword"
				},
				"content":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"reblogs_count":{
					"type":"long"
				},
				"favourites_count":{
					"type":"long"
				}
			}
		}
	}
}`

const WikiMapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"wiki":{
			"properties":{
				"id":{
					"type":"keyword"
				},
				"word":{
					"type":"keyword"
				},
				"created_at":{
					"type":"date"
				},
				"content":{
					"type":"text",
					"store": true,
					"fielddata": true
				}
			}
		}
	}
}`
