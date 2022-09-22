# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Bizflyctl < Formula
  desc " Command Line for Bizfly Cloud."
  homepage "https://github.com/bizflycloud/bizflyctl"
  version "0.2.6"

  on_macos do
    url "https://github.com/bizflycloud/bizflyctl/releases/download/v0.2.6/bizflyctl_Darwin_all.tar.gz"
    sha256 "52546ffa3c257471a4d332afe67380450dbc5755296eef1bf255ac04b8689370"

    def install
      bin.install "bizfly"
    end
  end

  on_linux do
    if Hardware::CPU.arm? && !Hardware::CPU.is_64_bit?
      url "https://github.com/bizflycloud/bizflyctl/releases/download/v0.2.6/bizflyctl_Linux_armv6.tar.gz"
      sha256 "972acb2c3b9a2fa5ea9ec2f32b525d69089dc801797e13902eb5b945e3710921"

      def install
        bin.install "bizfly"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/bizflycloud/bizflyctl/releases/download/v0.2.6/bizflyctl_Linux_x86_64.tar.gz"
      sha256 "7388a38d9b5fee810a3c27ae0038c2753dbcbb3d8a7c2d56a937d6c815e3f64e"

      def install
        bin.install "bizfly"
      end
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/bizflycloud/bizflyctl/releases/download/v0.2.6/bizflyctl_Linux_arm64.tar.gz"
      sha256 "4d9704517a5e08f6f9595dc1b649f1e6b5a776774d64a1c48661a0e58e5c8c45"

      def install
        bin.install "bizfly"
      end
    end
  end
end
