module Flame
  class Torrent < Model
    attribute :client, hidden: true
    attribute :infohash, convert: ->(value){value.downcase}
    attribute :name
    attribute :label
    attribute :total_peers
    attribute :peers
    attribute :total_seeds
    attribute :seeds
    attribute :eta
    attribute :queue
    attribute :size
    attribute :state
    attribute :path
    attribute :progress

    has_many :files
    child :stats, klass: Flame::Stats

    def add_file(file)
      files << file
    end

    def files_unsorted
      self.files
    end

    # def to_hash
    #   h = @data.dup
    #   h[:files] = files.map {|e| e.to_hash}
    #   h[:stats] = stats.to_hash
    #   h
    # end
  end
end
