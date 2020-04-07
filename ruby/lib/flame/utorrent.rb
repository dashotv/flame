module Flame
  class Utorrent < Client
    def list
      request("", {})
    end

    def add(uri)
      request("add", { url: uri })
    end

    %i{
      remove
      pause
      resume
      want_none
      stop
      start
    }.each do |n|
      define_method(n) do |infohash|
        request("#{n}", { infohash: infohash })
      end
    end

    def want(infohash, ids)
      request("want", { infohash: infohash, ids: ids })
    end

    def wanted?(infohash)
      request("wanted", { infohash: infohash })
    end
  end
end
