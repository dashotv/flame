RSpec.describe Flame do
  it "has a version number" do
    expect(Flame::Version::STRING).not_to be nil
  end
  let(:url) { "http://localhost:9001" }
  let(:client) { Flame::Client.new(url, {}) }

  context :nzbget do
    let(:nzb) { "https://api.nzbgeek.info/api?t=get&apikey=2b2f7303f77672ad619df8589e88b5d3&id=c4d0c7dee6bab6a35d9b74592ade0bb7" }

    it "can list" do
      expect { client.nzbget.list }.not_to raise_error
    end

    it "can add" do
      r = client.nzbget.add(nzb)
      expect(r["id"]).to be_an(Integer)
      expect(r["id"]).to be > 1
    end

    it "can pause" do
      sleep 1
      l = client.nzbget.list
      r = client.nzbget.pause(l["Result"].first['nzbid'])
      expect(r["error"]).to be false
    end

    it "can resume" do
      sleep 1
      l = client.nzbget.list
      r = client.nzbget.resume(l["Result"].first['nzbid'])
      expect(r["error"]).to be false
    end

    it "can remove" do
      sleep 1
      l = client.nzbget.list
      r = client.nzbget.remove(l["Result"].first['nzbid'])
      expect(r["error"]).to be false
    end

    it "can delete" do
      sleep 1
      l = client.nzbget.history(false)
      r = client.nzbget.remove(l["Result"].first['nzbid'], true)
      expect(r["error"]).to be false
    end

    it "can destroy" do
      sleep 1
      l = client.nzbget.history(true)
      r = client.nzbget.destroy(l["Result"].first['nzbid'])
      expect(r["error"]).to be false
    end
  end

  context :utorrent do
    let(:hash) { "6f8cd699135b491513e65d967a052a7087750d9c" }
    let(:torrent) { "http://www.slackware.com/torrents/slackware-14.1-install-d1.torrent" }

    it "can list" do
      expect { client.utorrent.list }.not_to raise_error
    end

    it "can add" do
      r = client.utorrent.add(torrent)
      expect(r["infohash"]).to eq(hash)
    end

    it "can pause" do
      sleep 3
      l = client.utorrent.list
      r = client.utorrent.pause(l["Torrents"].first["Hash"])
      expect(r["error"]).to be false
    end

    it "can resume" do
      l = client.utorrent.list
      r = client.utorrent.resume(l["Torrents"].first["Hash"])
      expect(r["error"]).to be false
    end

    it "can start" do
      l = client.utorrent.list
      r = client.utorrent.start(l["Torrents"].first["Hash"])
      expect(r["error"]).to be false
    end

    it "can stop" do
      l = client.utorrent.list
      r = client.utorrent.stop(l["Torrents"].first["Hash"])
      expect(r["error"]).to be false
    end

    it "can remove" do
      l = client.utorrent.list
      r = client.utorrent.remove(l["Torrents"].first["Hash"])
      expect(r["error"]).to be false
    end
  end

  context :qbittorrent do
    let(:hash) { "6f8cd699135b491513e65d967a052a7087750d9c" }
    let(:torrent) { "http://www.slackware.com/torrents/slackware-14.1-install-d1.torrent" }

    it "can list" do
      expect { client.qbittorrent.list }.not_to raise_error
    end

    it "can add" do
      r = client.qbittorrent.add(torrent)
      expect(r["infohash"]).to eq(hash)
    end

    it "can pause" do
      sleep 3
      l = client.qbittorrent.list
      r = client.qbittorrent.pause(l["Torrents"].first["Hash"])
      expect(r["error"]).to be false
    end

    it "can resume" do
      l = client.qbittorrent.list
      r = client.qbittorrent.resume(l["Torrents"].first["Hash"])
      expect(r["error"]).to be false
    end

    it "can remove" do
      l = client.qbittorrent.list
      r = client.qbittorrent.remove(l["Torrents"].first["Hash"])
      expect(r["error"]).to be false
    end
  end
end
