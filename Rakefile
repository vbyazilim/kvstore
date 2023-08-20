namespace :run do
  desc "run server"
  task :server do
    system %{
      SERVER_ENV=#{ENV['SERVER_ENV'] || "local"} go run -race cmd/server/main.go
    }
  end
end

desc "default task"
task :default => ["run:server"]