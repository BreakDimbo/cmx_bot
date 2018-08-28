#!/Users/break/.rvm/rubies/ruby-2.4.2/bin/ruby
require 'json'
require 'rest-client'

LOCAL_URL = "http://localhost:9200"
REMOTE_URL = "http://47.93.43.59:9201"
USERNAME = "elastic"
PASSWORD = "break12345"


def update_mapping_add_field(field_name, default_value, index, type)
  query = "/#{index}/_mapping/#{type}"
end

def set_alias(production_index, current_index)
  query = "/_aliases"
  payload = {"actions": [
    {"add": {"index": "#{production_index}", "alias": "#{current_index}"}}
  ]}
  p request(payload, query, :post)
end

def create_new_index(new_index, old_index, type, field_name, field_type)
  query = "/#{new_index}"
  mapping = get_current_index_mapping(old_index)
  mapping["#{type}"]["properties"]["#{field_name}"] = {"type": "#{field_type}"}
  payload = {
    "settings": {
    "number_of_shards": 1,
		"number_of_replicas": 0
    }, 
    "mappings": mapping
  }
  
  p request(payload, query, :put)
end

def copy_data(from, to)
  query = "/_reindex"
  payload = {
    "source": {
      "index": from
    },
    "dest": {
      "index": to
    }
  }
  p request(payload, query, :post)
end

def set_default(index, type, field_name, default_value)
  query = "/#{index}/_update_by_query"
  payload = {
    "script":{
      "source": "ctx._source.#{field_name}=\"#{default_value}\""
    }
  }
  p request(payload, query, :post)
end

def change_alias(new_index, aliass)
  query = "/_aliases"
  payload = {"actions": [
    # {"remove": {"index": "#{old_index}", "alias": "#{aliass}"}},
    {"add": {"index": "#{new_index}", "alias": "#{aliass}"}}
  ]}
  p request(payload, query, :post)
end

def delete(index)
  query = "/#{index}"
  p request(nil, query, :delete)
end

def get_current_index_mapping(index)
  query = "/#{index}"
  response = request(nil, query, :get)
  return unless response[0].equal?(:success)
  mapping = response[1]["#{index}"]["mappings"]
end

def parse_json(response)
  JSON.parse(response)
end

def request(payload, query, method)
  response = RestClient::Request.new({
    method: method,
    url: "#{REMOTE_URL}#{query}",
    user: USERNAME,
    password: PASSWORD,
    payload: payload.to_json,
    headers: { :accept => :json, content_type: :json }
  }).execute do |response, request, result|
    case response.code
    when 400
      [ :error, parse_json(response.to_str) ]
    when 200
      [ :success, parse_json(response.to_str) ]
    else
      fail "Invalid response #{response.to_str} received."
    end
  end
end

# set_alias("status", "status_v1")
create_new_index("local_v2", "local", "status", "server", "keyword")
copy_data("local", "local_v2")
set_default("local","status","server","https://cmx.im")
delete("local")
change_alias("local_v2", "local")