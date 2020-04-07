module Flame
  class Utorrent < Client
    def list
      response = request('/', {})
      puts "response: #{response.inspect}"
    end

    def add

    end
  end
end
