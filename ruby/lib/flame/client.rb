require 'rest-client'
require 'json'

module Flame
  class Client
    def initialize(url, options={})
      @base = url
      @options = {}.merge(options)
      @headers = {}
    end

    def nzbget
      @nzbget ||= Flame::Nzbget.new("#{@base}/nzbs", @options)
    end

    def utorrent
      @utorrent ||= Flame::Utorrent.new("#{@base}/torrents", @options)
    end

    private

    def request(url, data)
      query = query(data)
      opt = {
        :method => :get,
        :url => "#{@base}/#{url}?#{query}",
        :headers => @headers,
        :verify_ssl => false,
      }
      json = RestClient::Request.execute(opt)
      puts "json: #{json}"
      raise :empty_response unless json && json != ''
      ::JSON.parse(json)
    end

    def query(hash)
      hash.map { |k, v| "#{k}=#{v}" }.join('&')
    end
  end
end
