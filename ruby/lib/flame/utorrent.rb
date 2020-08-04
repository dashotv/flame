require 'ostruct'

module Flame
  class Utorrent < Client
    def list
      data = request("", {})
      response(data)
    end

    def add(uri)
      u = Base64.encode64(uri).chomp
      request("add", { url: u })
    end

    %i{
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

    def remove(infohash, delete = false)
      request("remove", { infohash: infohash, delete: delete })
    end

    def want(infohash, ids)
      request("want", { infohash: infohash, ids: ids })
    end

    def wanted?(infohash)
      request("wanted", { infohash: infohash })
    end

    protected

    def response(data)
      r    = Flame::Response.new(self)
      up   = 0
      down = 0
      data['Torrents'].each do |d|
        t = OpenStruct.new(d)

        torrent             = Flame::Torrent.new
        torrent.client      = self
        torrent.infohash    = t.Hash
        torrent.name        = t.Name
        torrent.label       = t.Label
        torrent.progress    = t.Progress
        torrent.seeds       = t.SeedsConnected
        torrent.total_seeds = t.SeedsTotal
        torrent.peers       = t.PeersConnected
        torrent.total_peers = t.PeersTotal
        torrent.eta         = t.Finish
        torrent.queue       = t.Queue
        torrent.size        = t.Size
        torrent.state       = t.State
        torrent.path        = t.Path

        up   += t.UploadRate
        down += t.DownloadRate

        torrent.stats.load({
            upload:   { rate: t.UploadRate / 1000.0 },
            download: { rate: t.DownloadRate / 1000.0 }
        })

        t.Files.each do |df|
          f = OpenStruct.new(df)

          file          = Flame::File.new(name: "#{torrent.path}/#{f.Name}", num: f.Number, size: f.Size, priority: f.Priority)
          file.progress = (f.Downloaded / f.Size.to_f) * 100
          torrent.add_file(file)
        end

        r.add(torrent)
      end

      r.stats.load({ upload: { rate: "%.2f" % [up / 1024.0] }, download: { rate: "%.2f" % [down / 1024.0] } })
      r
    end
  end
end
