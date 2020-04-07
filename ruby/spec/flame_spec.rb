RSpec.describe Flame do
  it "has a version number" do
    expect(Flame::VERSION).not_to be nil
  end

  context :nzbget do
    let(:client) { Flame::Client.new("http://localhost:3001", {}) }
    it "can list" do
      expect { client.nzbget.list }.not_to raise_error
    end
  end

  context :utorrent do
    let(:client) { Flame::Client.new("http://localhost:3001", {}) }
    it "can list" do
      expect { client.utorrent.list }.not_to raise_error
    end
  end
end
