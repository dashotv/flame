RSpec.describe Flame do
  it "has a version number" do
    expect(Flame::VERSION).not_to be nil
  end

  context :nzbget do
    let(:client) { Flame::Client.new("http://localhost:3001", {}) }
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
    let(:client) { Flame::Client.new("http://localhost:3001", {}) }
    # it "can list" do
    #   expect { client.utorrent.list }.not_to raise_error
    # end
  end
end
