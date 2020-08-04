module Flame
  class Model
    def initialize(data={})
      @data = {}.merge(data)
      self.class.attributes.each do |name, options|
        self.send(name)
      end
    end

    def [](key)
      self.send(key)
    end

    def []=(key, value)
      self.send("#{key}=", value)
    end

    def to_hash
      out = {}
      self.class.attributes.each do |name, options|
        next if options[:hidden]
        type = options[:type]
        # puts "to_hash:#{self.class.name}:#{name}:#{type}"
        case type
          when :string
            out[name] = @data[name]
          when :array
            out[name] = @data[name].map {|e| e.to_hash}
          when :object
            out[name] = @data[name].to_hash
          else
            out[name] = @data[name]
        end
      end
      out#.delete_if {|k, v| v.nil?}
    end

    class << self
      attr_reader :attributes

      protected

      def attribute(name, o={})
        namesym = name.to_sym
        options = {type: :string}.merge(o)

        @attributes ||= {}
        @attributes[namesym] = options

        default = options[:default] || nil
        define_method "#{name}" do
          # puts "get:#{name}"
          @data[namesym] ||= default.is_a?(Proc) ? default.call : default
          @data[namesym]
        end
        define_method "#{name}=" do |value|
          @data[namesym] ||= default.is_a?(Proc) ? default.call : default
          # old = @data[namesym]
          new = value
          if options[:convert]
            new = options[:convert].call(value)
          end
          @data[namesym] = new
          # puts "#{name}= #{old} => #{new}"
          new
        end
      end
      def has_many(name, options={})
        attribute(name, {default: ->{[]}, type: :array}.merge(options))
      end
      def child(name, options={})
        klass = options.delete(:klass)
        attribute(name, {default: ->{klass.new}, type: :object}.merge(options))
      end
    end
  end
end
