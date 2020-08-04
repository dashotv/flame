require 'ostruct'

module Flame
  class Qbittorrent < Client
    def list
      data = request("", {})
      response(data)
    end

    def add(uri)
      u = Base64.encode64(uri).chomp
      request("add", { url: u })
    end

    %i{
      remove
      pause
      resume
      want_none
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

    protected

    def response(data)
      r = Flame::Response.new(self)
      up = data["UploadRate"].to_i
      down = data["DownloadRate"].to_i
      data['Torrents'].each do |d|
        t = OpenStruct.new(d)

        torrent = Flame::Torrent.new
        torrent.client = self
        torrent.infohash = t.Hash
        torrent.name = t.Name
        torrent.label = t.Label
        torrent.progress = t.Progress
        torrent.eta = t.Finish
        torrent.queue = t.Queue
        torrent.size = t.Size
        torrent.state = t.State
        torrent.path = t.Path

        torrent.stats.load({
            upload: {rate: t.UploadRate/1000.0},
            download: {rate: t.DownloadRate/1000.0}
        })

        t.Files.each do |df|
          f = OpenStruct.new(df)

          file = Flame::File.new(name: "#{torrent.path}/#{f.name}", num: f.id, size: f.size, priority: f.priority)
          file.progress = f.progress
          torrent.add_file(file)
        end

        r.add(torrent)
      end

      r.stats.load({upload: {rate: "%.2f" % [up/1024.0]}, download: {rate: "%.2f" % [down/1024.0]}})
      r
    end
  end
end
