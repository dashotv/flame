module Flame
  class Response
    attr_reader :stats

    def initialize(client=nil)
      @client = client
      @torrents = {}
      @count = 0
      @stats = Flame::Stats.new
    end

    def [](infohash)
      get(infohash)
    end

    def has?(infohash)
      @torrents.keys.include?(infohash.downcase)
    end

    def get(infohash)
      raise "torrent not found: #{infohash} :: #{@torrents.keys.inspect}" unless @torrents[infohash.downcase]
      @torrents[infohash.downcase]
    end

    def add(torrent)
      @torrents[torrent.infohash.downcase] = torrent
      @count += 1
    end

    def remove(infohash)
      torrent = @torrents.delete(infohash)
      @count -= 1 unless torrent.nil?
      torrent
    end

    def hashes
      @torrents.keys
    end

    def each(&block)
      if block_given?
        sorted_torrents.each do |k, v|
          case block.arity
            when 1
              yield v
            when 2
              yield k, v
            else
              raise "unknown arity for block in each: #{block.arity}"
          end
        end
      end
    end

    def to_hash
      torrents = sorted_torrents.inject([]) {|a, e| a << e.last.to_hash}
      {
          torrents: torrents,
          count: @count,
          stats: @stats.to_hash
      }
    end

    private

    def sorted_torrents
      @torrents.sort_by{|k, v| v.queue}
    end
  end
end
