RSpec.describe Flame::Utorrent do
  let(:url) { "http://localhost:9001" }
  let(:client) { Flame::Client.new(url, {}) }
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
    r = client.utorrent.pause(hash)
    expect(r["error"]).to be false
  end

  it "can resume" do
    r = client.utorrent.resume(hash)
    expect(r["error"]).to be false
  end

  it "can start" do
    r = client.utorrent.start(hash)
    expect(r["error"]).to be false
  end

  it "can stop" do
    r = client.utorrent.stop(hash)
    expect(r["error"]).to be false
  end

  it "can remove" do
    r = client.utorrent.remove(hash)
    expect(r["error"]).to be false
  end
end
