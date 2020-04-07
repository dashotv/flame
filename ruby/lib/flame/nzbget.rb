module Flame
  class Nzbget < Client
    def list
      response = request('/', {})
      puts "response: #{response.inspect}"
    end

    def add

    end
  end
end
