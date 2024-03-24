/** @type {import('next').NextConfig} */
const nextConfig = {
    webpack: (config, { isServer }) => {
      if (!isServer) {
        config.resolve.fallback = {
          ...config.resolve.fallback,
          child_process: false,
          crypto: false,
          'fs/promises': false,
          net: false,
          tls:false,
          fs:false,
          dns:false, 
          kerberos: false,
        '@mongodb-js/zstd': false,
        '@aws-sdk/credential-providers': false,
        'gcp-metadata': false,
        snappy: false,
        socks: false,
        "aws4": false,
        'mongodb':false,
        'mongodb-client-encryption': false,
        };
      }
      return config;
    },
  };
  
  export default nextConfig;