task :command_exists, [:command] do |_, args|
  abort "#{args.command} doesn't exists" if `command -v #{args.command} > /dev/null 2>&1 && echo $?`.chomp.empty?
end

task :has_bumpversion do
  Rake::Task['command_exists'].invoke('bumpversion')
end

task :is_repo_clean do
  abort 'please commit your changes first!' unless `git status -s | wc -l`.strip.to_i.zero?
end

AVAILABLE_REVISIONS = %w[major minor patch].freeze
task :bump, [:revision] => [:has_bumpversion] do |_, args|
  args.with_defaults(revision: 'patch')
  unless AVAILABLE_REVISIONS.include?(args.revision)
    abort "Please provide valid revision: #{AVAILABLE_REVISIONS.join(',')}"
  end

  system "bumpversion #{args.revision}"
  exit $?.exitstatus
end

namespace :run do
  desc "run server"
  task :server do
    system %{
      source .env &&
      go run -race cmd/server/main.go
    }
    exit $?.exitstatus
  end
end

desc "default task"
task :default => ["run:server"]

desc "run golangci lint"
task :lint do
  system %{
    LOG_LEVEL=error golangci-lint run
  }
  exit $?.exitstatus
end

DOCKER_IMAGE_TAG = 'kvstore:latest'
DOCKER_ENV_VARS = %w[
  SERVER_ENV
]
DOCKER_HOST_PORT = ENV['DOCKER_HOST_PORT'] || 9000
namespace :docker do
  desc "Build image (locally)"
  task :build do
    git_commit_hash = `git rev-parse HEAD`.chomp
    goos =`go env GOOS`.chomp
    goarch =`go env GOARCH`.chomp
    build_commit_hash = "#{git_commit_hash}-#{goos}-#{goarch}"
    
    system %{
      docker build --build-arg="BUILD_INFORMATION=#{build_commit_hash}" \
        -t #{DOCKER_IMAGE_TAG} .
    }
    exit $?.exitstatus
  end
  # -p HOST:CONTAINER
  desc "Run image (locally)"
  task :run do
    system %{
      source .env &&
      echo "service will be available on port: #{DOCKER_HOST_PORT}" &&
      echo &&
      docker run \
        --cpus="2" \
        --env #{DOCKER_ENV_VARS.join(" --env ")} \
        -p #{DOCKER_HOST_PORT}:8000 \
        #{DOCKER_IMAGE_TAG}
    }
    exit $?.exitstatus
  end
end

desc "release new version #{AVAILABLE_REVISIONS.join(',')}, default: patch"
task :release, [:revision] => [:is_repo_clean] do |_, args|
  args.with_defaults(revision: 'patch')
  Rake::Task['bump'].invoke(args.revision)
end

namespace :test do
  desc "run all tests"
  task :run_all do
    system %{
      LOG_LEVEL="error" go test -race -p 1 -v -race -failfast ./...
    }
  end
end
