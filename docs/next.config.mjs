import nextMDX from '@next/mdx'
import { createContentlayerPlugin } from 'next-contentlayer2'
import { mdxOptions } from './src/mdx/mdx-options.mjs'
import path from 'path'

const withMDX = nextMDX({
  options: mdxOptions,
})

/** @type {import('next').NextConfig} */
const nextConfig = {
  basePath: '/docs',
  pageExtensions: ['js', 'jsx', 'ts', 'tsx', 'mdx'],
  experimental: {
    outputFileTracingIncludes: {
      '/**/*': ['./src/app/**/*.mdx'],
    },
    serverActions: {
      allowedOrigins: ['nitric.io'],
    },
    scrollRestoration: true,
  },
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'github.com',
        port: '',
      },
      {
        protocol: 'https',
        hostname: 'raw.githubusercontent.com',
        port: '',
      },
    ],
  },
  async redirects() {
    return [
      {
        source: '/',
        destination: '/docs',
        basePath: false,
        permanent: false,
      },
      {
        source: '/docs/reference',
        destination: '/docs',
        basePath: false,
        permanent: true,
      },
      // redirects from old docs
      ...[
        '/docs/reference/go/secrets/secret-version',
        '/docs/reference/go/secrets/secret-latest',
        '/docs/reference/go/secrets/secret-version-access',
        '/docs/reference/go/storage/bucket-file',
        '/docs/reference/go/storage/bucket-files',
        '/docs/reference/go/storage/bucket-file-read',
        '/docs/reference/go/storage/bucket-file-write',
        '/docs/reference/go/storage/bucket-file-delete',
        '/docs/reference/go/storage/bucket-file-downloadurl',
        '/docs/reference/go/storage/bucket-file-uploadurl',
      ].map((source) => ({
        source: source,
        destination: `/docs/reference/go`,
        basePath: false,
        permanent: true,
      })),
      ...[
        '/docs/testing',
        '/docs/guides/debugging',
        '/docs/guides/serverless-rest-api-example',
        '/docs/guides/graphql',
        '/docs/guides/serverless-api-with-planetscale-and-prisma',
        '/docs/guides/nitric-and-supabase',
        '/docs/guides/api-with-nextjs',
        '/docs/guides/twilio',
        '/docs/guides/stripe',
        '/docs/guides/secure-api-auth0',
        '/docs/guides/byo-database',
      ].map((source) => ({
        source: source,
        destination: `/docs/guides/nodejs${source.replace(
          /^(\/docs\/guides\/|\/docs\/)/,
          '/',
        )}`,
        basePath: false,
        permanent: true,
      })),
      ...['/docs/guides/text-prediction', '/docs/guides/create-histogram'].map(
        (source) => ({
          source: source,
          destination: `/docs/guides/python${source.replace(
            /^(\/docs\/guides\/|\/docs\/)/,
            '/',
          )}`,
          basePath: false,
          permanent: true,
        }),
      ),
      {
        source: '/docs/comparison/:slug',
        destination: '/docs/concepts/comparison/:slug',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/assets/faq',
        destination: '/docs/faq',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/assets',
        destination: '/docs',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/assets/eject',
        destination: '/docs/faq',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/concepts',
        destination: '/docs',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/env',
        destination: '/docs/reference/env',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/assets/custom-containers',
        destination: '/docs/reference/custom-containers',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/support/eject',
        destination: '/docs/concepts/eject',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/access-control',
        destination: '/docs/get-started/foundations/infrastructure/security',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/assets/comparison/:path*',
        destination: '/docs/faq/comparison/:path*',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/getting-started',
        destination: '/docs/get-started/quickstart',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/guides/getting-started',
        destination: '/docs/get-started/quickstart',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/guides/getting-started/concepts',
        destination: '/docs',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/getting-started/concepts',
        destination: '/docs',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/assets/env',
        destination: '/docs/reference/env',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/assets/resources-overview',
        destination: '/docs/getting-started/resources-overview',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/assets/examples',
        destination: '/docs/guides/examples',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/assets/access-control',
        destination: '/docs/get-started/foundations/infrastructure/security',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/concepts/eject',
        destination: '/docs/faq',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/guides/getting-started/deploy',
        destination: '/docs/getting-started/deployment',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/guides/getting-started/language-support',
        destination: '/docs/reference/languages#supported-languages',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/guides/getting-started/nodejs/:slug*',
        destination: '/docs/guides/nodejs/:slug*',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/guides/getting-started/dart/:slug*',
        destination: '/docs/guides/dart/:slug*',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/guides/getting-started/go/:slug*',
        destination: '/docs/guides/go/:slug*',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/guides/getting-started/jvm/:slug*',
        destination: '/docs/guides/jvm/:slug*',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/guides/getting-started/python/:slug*',
        destination: '/docs/guides/python/:slug*',
        basePath: false,
        permanent: true,
      },
      ...[
        'azure-pipelines',
        'gitlab-ci',
        'github-actions',
        'google-cloud-build',
      ].map((source) => ({
        source: `/docs/guides/getting-started/${source}`,
        destination: `/docs/guides/deploying/${source}`,
        basePath: false,
        permanent: true,
      })),
      {
        source: '/docs/guides/getting-started/quickstart',
        destination: '/docs/get-started/quickstart',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/installation',
        destination: '/docs/get-started/installation',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/guides/getting-started/installation',
        destination: '/docs/get-started/installation',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/guides/getting-started/local-dashboard',
        destination: '/docs/get-started/foundations/projects/local-development',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/guides/getting-started/deployment',
        destination: '/docs/get-started/foundations/deployment',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/faq/common-questions',
        destination: '/docs/faq',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/pulumi',
        destination: '/docs/providers/pulumi',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/pulumi/custom-providers',
        destination: '/docs/providers/custom',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/pulumi/pulumi-cloud',
        destination:
          '/docs/providers/pulumi#using-your-pulumi-cloud-account-optional',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/dart/v1',
        destination: '/docs/reference/dart',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/go/api/api-details',
        destination: '/docs/reference/go',
        basePath: false,
        permanent: true,
      },
      ...['nodejs', 'dart', 'go', 'jvm', 'python', 'terraform'].map((page) => ({
        source: `/docs/guides/${page}`,
        destination: `/docs/guides`,
        basePath: false,
        permanent: true,
      })),
      {
        source: '/docs/concepts/access-control',
        destination: '/docs/get-started/foundations/infrastructure/security',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/concepts/build-with-nitric',
        destination: '/docs/get-started/foundations/why-nitric',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/concepts/cicd',
        destination: '/docs/get-started/foundations/deployment#ci-cd',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/concepts/comparison',
        destination: '/docs/misc/faq#differences-from-other-solutions',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/concepts/comparison/:slug*',
        destination: '/docs/misc/comparison/:slug*',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/concepts/extensibility',
        destination: '/docs/get-started/foundations/deployment#flexibility',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/concepts/how-devs-use-nitric',
        destination: '/docs/get-started/foundations/why-nitric',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/concepts/how-nitric-works',
        destination: '/docs/get-started/foundations/deployment',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/concepts/how-ops-use-nitric',
        destination: '/docs/get-started/foundations/deployment',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/concepts/infrastructure-from-code',
        destination: '/docs/get-started/foundations/infrastructure',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/concepts/introduction',
        destination: '/docs',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/concepts/language-support',
        destination: '/docs/reference/languages#supported-languages',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/concepts/project-structure',
        destination: '/docs/get-started/foundations/projects#project-structure',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/cli/installation',
        destination: '/docs/get-started/installation',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/cli/local-development',
        destination: '/docs/get-started/foundations/projects/local-development',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/cli/stacks',
        destination: '/docs/get-started/foundations/deployment#stacks',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/guides/examples',
        destination: 'https://github.com/nitrictech/examples',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/guides/deploying',
        destination:
          '/docs/get-started/foundations/deployment#deploying-your-application',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/getting-started/deployment',
        destination: '/docs/get-started/foundations/deployment',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/getting-started/installation',
        destination: '/docs/get-started/installation',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/getting-started/local-dashboard',
        destination: '/docs/get-started/foundations/projects/local-development',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/getting-started/quickstart',
        destination: '/docs/get-started/quickstart',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/getting-started/resources-overview',
        destination: '/docs/get-started/foundations/infrastructure/resources',
        basePath: false,
        permanent: true,
      },
      ...[
        'apis',
        'keyvalue',
        'queues',
        'topics',
        'schedules',
        'secrets',
        'storage',
      ].flatMap((page) => [
        {
          source: `/docs/providers/pulumi/aws/${page}`,
          destination: `/docs/providers/mappings/aws/${page}`,
          basePath: false,
          permanent: true,
        },
        {
          source: `/docs/providers/pulumi/azure/${page}`,
          destination: `/docs/providers/mappings/azure/${page}`,
          basePath: false,
          permanent: true,
        },
        {
          source: `/docs/providers/pulumi/gcp/${page}`,
          destination: `/docs/providers/mappings/gcp/${page}`,
          basePath: false,
          permanent: true,
        },
      ]),
      {
        source: '/docs/reference/providers',
        destination: '/docs/providers',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/providers/aws',
        destination: '/docs/providers/pulumi/aws',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/providers/pulumi/aws/imports',
        destination: '/docs/providers/pulumi/aws#importing-existing-resources',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/providers/aws/configuration',
        destination: '/docs/providers/pulumi/aws#stack-configuration',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/providers/pulumi/aws/configuration',
        destination: '/docs/providers/pulumi/aws#stack-configuration',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/providers/aws/:slug((?!configuration).*)',
        destination: '/docs/providers/pulumi/aws/:slug*',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/providers/gcp',
        destination: '/docs/providers/pulumi/gcp',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/providers/gcp/configuration',
        destination: '/docs/providers/pulumi/gcp#stack-configuration',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/providers/gcp/:slug((?!configuration).*)',
        destination: '/docs/providers/pulumi/gcp/:slug*',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/providers/azure',
        destination: '/docs/providers/pulumi/azure',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/providers/azure/configuration',
        destination: '/docs/providers/pulumi/azure#stack-configuration',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/providers/azure/:slug((?!configuration).*)',
        destination: '/docs/providers/pulumi/azure/:slug*',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/providers/custom/building-custom-provider',
        destination: '/docs/providers/custom/create',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/providers/custom/extend-standard-provider',
        destination: '/docs/providers/custom/extend',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/providers/pulumi',
        destination: '/docs/providers/pulumi',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/providers/pulumi/custom',
        destination: '/docs/providers/custom',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/providers/terraform',
        destination: '/docs/providers/terraform',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/providers/terraform/:slug*',
        destination: '/docs/providers/terraform/:slug*',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/reference/providers/install/docker',
        destination: '/docs/providers/custom/docker',
        basePath: false,
        permanent: true,
      },

      ...['contributions', 'faq', 'support', 'support/upgrade'].map((page) => ({
        source: `/docs/${page}`,
        destination: `/docs/misc/${page}`,
        basePath: false,
        permanent: true,
      })),
      {
        source: '/docs/reference/python/schedules/schedule',
        destination: '/docs/reference/python/schedule/schedule',
        basePath: false,
        permanent: true,
      },
      {
        source: '/docs/misc/comparison/winglang',
        destination: '/docs/',
        basePath: false,
        permanent: true,
      },
    ]
  },
  async headers() {
    return [
      {
        source: '/:path*',
        headers: [
          {
            key: 'Strict-Transport-Security',
            value: '',
          },
          {
            key: 'X-Robots-Tag',
            value: 'all',
          },
          {
            key: 'X-Frame-Options',
            value: 'DENY',
          },
        ],
      },
    ]
  },
}

const withContentlayer = createContentlayerPlugin({
  configPath: path.resolve(process.cwd(), './contentlayer.config.ts'),
})

export default withContentlayer(withMDX(nextConfig))
