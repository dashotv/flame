module Flame
  class Stats < Model

    # class Rate < Model
    #   attribute :rate
    #   attribute :limit
    #   attribute :overhead
    # end
    class Upload < Model
      attribute :rate
      attribute :limit
      attribute :overhead
    end
    class Download < Model
      attribute :rate
      attribute :limit
      attribute :overhead
    end
    class Connection < Model
      attribute :incoming
      attribute :count
      attribute :limit
    end
    class Space < Model
      attribute :free
    end
    class Dht < Model
      attribute :nodes
    end

    child :upload, klass: Flame::Stats::Upload
    child :download, klass: Flame::Stats::Download
    child :connection, klass: Flame::Stats::Connection
    child :dht, klass: Flame::Stats::Dht
    child :space, klass: Flame::Stats::Space

    def load(data)
      #puts "load_stats"
      #puts "data: #{data.inspect}"
      data.each do |key, values|
        os = self.send(key)
        #puts "#{key} = #{os.inspect}"
        values.each do |vk, val|
          #puts "  #{vk} = #{val}"
          os.send("#{vk}=", val)
        end
        #puts "#{key} = #{os.inspect}"
      end
      #puts "stats load: #{self.inspect}"
    end
  end
end
