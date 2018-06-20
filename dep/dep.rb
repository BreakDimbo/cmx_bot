#!/Users/break/.rvm/rubies/ruby-2.4.2/bin/ruby

# usage ./dep.rb -s[service] -t[target] -d[work_directory]
# example dep.rb -s bot,fbot -t break.bwh.com -d /Users/break/Work/Workspace/udesk/udesk_qilin_cti
# use default config: dep.rb -s bot,fbot -t break.bwh.com
require 'optparse'

DefaultDir = "/Users/break/Documents/Geek/cmx_bot"

options = {}
OptionParser.new do |opts|
  opts.banner = "Usage: ./dep.rb -s[services] -t[targets] -d[work_directory]"

  opts.on("-s", "--services [Services]", Array, "make services") do |v|
    options[:services] = v
  end

  opts.on("-t", "--targets [Target]", Array, "target servers") do |v|
    options[:targets] = v
  end

  opts.on("-d", "--dir [Dir]", String, "directory of project makefile") do |v|
    options[:dir] = v
  end
end.parse!

p options

def pull_and_compile_restart(services, target, dir)
  services.each do |service|
    puts "start pull #{service} to #{target}"
    mkdir = "intelbot" if service == "bot"
    mkdir = "firebot" if service == "fbot"

    raise unless system("ssh #{target} \"cd /home/break/documents/cmx_bot; export GOPATH=/home/break/documents/cmx_bot; git pull; /usr/local/go/bin/go build -o bin/#{service} bot/#{mkdir}\"")
    raise unless system("ssh #{target} \"sudo systemctl stop #{service}.service\"")
    raise unless system("ssh #{target} \"cp /home/break/documents/cmx_bot/bin/#{service} /usr/local/cmx_bot/current/bin\"")
    raise unless system("ssh #{target} \"sudo systemctl start #{service}.service\"")
    puts "over restart #{service} on #{target}"
  end
end

def check
  # TODO 检查版本
end

def execute(options)
  work_dir = options[:dir] || DefaultDir
  services = options[:services] || ["bot"]
  targets = options[:targets] || ["break.bwh.com"]
  Dir.chdir(work_dir)
  
  targets.each do |target|
    pull_and_compile_restart(services,target,work_dir)
    check()
  end
end

execute(options)