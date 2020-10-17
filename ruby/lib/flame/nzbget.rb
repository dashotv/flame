require 'base64'

module Flame
  class Nzbget < Client
    def list
      request("", {})
    end

    def add(url, options={})
      u = Base64.encode64(url).chomp.gsub(/[\r\n]+/, '')
      request("add", { url: u }.merge(options))
    end

    %i{
      pause
      resume
      destroy
    }.each do |n|
      define_method(n) do |id|
        request("#{n}", { id: id })
      end
    end

    def remove(id, delete=false)
      request("remove", {id: id, delete: delete})
    end

    def history(hidden=false)
      request("history", {hidden: hidden})
    end
  end
end
