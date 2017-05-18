require 'rake'
require 'bundler/setup'
require 'haml'

task :haml => %W[index.tmpl]

task :test do |t, args|
  puts  "test"
end

rule ".tmpl" => ".haml" do |t|
  haml_template = Haml::Engine.new File.read(t.source), {ugly: false}
  File.write t.name, haml_template.render
  
  puts "Compiled #{t.name}"
end