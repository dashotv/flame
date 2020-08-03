require 'rest-client'
require 'json'

module Flame
  class EmptyResponseError < StandardError; end
  class BadRequestError < StandardError; end
  class Client
    def initialize(url, options = {})
      @base    = url
      @options = {}.merge(options)
      @headers = {}
    end

    def nzbget
      @nzbget ||= Flame::Nzbget.new("#{@base}/nzbs", @options)
    end

    def utorrent
      @utorrent ||= Flame::Utorrent.new("#{@base}/torrents", @options)
    end

    def qbittorrent
      @utorrent ||= Flame::Qbittorrent.new("#{@base}/qbittorrents", @options)
    end

    private

    def request(path, params)
      url  = "#{@base}/#{path}?#{params.to_query}"
      opt  = {
        method:     :get,
        url:        url,
        headers:    @headers,
        verify_ssl: false,
      }
      err = nil
      begin
        json = RestClient::Request.execute(opt)
      rescue RestClient::ExceptionWithResponse => e
        err = e
        json = e.response
      end
      raise EmptyResponseError.new("empty response") unless json && json != ''
      parsed = ::JSON.parse(json)
      raise BadRequestError.new("#{err.message}: #{parsed["error"]}") if err
      parsed
    end
  end
end
