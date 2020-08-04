require 'rest-client'
require 'uri'
require 'json'

module Flame
  class EmptyResponseError < StandardError; end

  class BadRequestError < StandardError; end

  class Client
    def initialize(url, options = {})
      @base    = url
      @options = {}.merge(options)
      @headers = {
          content_type: 'application/json',
          accept:       'application/json',
      }
    end

    def nzbget
      @nzbget ||= Flame::Nzbget.new("#{@base}/nzbs", @options)
    end

    def utorrent
      @utorrent ||= Flame::Utorrent.new("#{@base}/torrents", @options)
    end

    def qbittorrent
      @qbittorrent ||= Flame::Qbittorrent.new("#{@base}/qbittorrents", @options)
    end

    private

    def request(path, params)
      url = "#{@base}/#{path}?#{URI.encode_www_form(params)}"
      opt = {
          method:     :get,
          url:        url,
          headers:    @headers,
          verify_ssl: false,
      }
      # puts "request: #{url}"
      err = nil
      begin
        json = RestClient::Request.execute(opt)
      rescue RestClient::ExceptionWithResponse => e
        err  = e
        json = e.response
      end
      raise EmptyResponseError.new("empty response") unless json && json != ''
      parsed = ::JSON.parse(json)
      raise BadRequestError.new("#{err.message}: #{parsed["error"]}") if err
      parsed
    end
  end
end
