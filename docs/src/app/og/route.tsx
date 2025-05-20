import { BASE_URL } from '@/lib/constants'
import { ImageResponse } from 'next/og'
import { NextRequest } from 'next/server'

export const runtime = 'edge'

const size = {
  width: 1200,
  height: 630,
}

const font = fetch(new URL('../../assets/Sora-Bold.ttf', import.meta.url)).then(
  (res) => res.arrayBuffer(),
)

// Image generation
export async function GET(req: NextRequest) {
  const { searchParams } = req.nextUrl

  const title = searchParams.get('title')?.slice(0, 100) || 'Nitric Docs'

  const description =
    searchParams.get('description')?.slice(0, 100) ||
    'Docs for the Nitric cloud application framework.'

  return new ImageResponse(
    (
      // ImageResponse JSX element
      <div
        style={{
          backgroundColor: '#121118',
          backgroundSize: '150px 150px',
          height: '100%',
          width: '100%',
          display: 'flex',
          textAlign: 'center',
          alignItems: 'center',
          justifyContent: 'center',
          flexDirection: 'column',
          flexWrap: 'nowrap',
          position: 'relative',
        }}
      >
        <div
          style={{
            display: 'flex',
            alignItems: 'center',
            alignContent: 'center',
            justifyContent: 'center',
            justifyItems: 'center',
          }}
        >
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img alt="Nitric" src={`${BASE_URL}/docs/images/open_graph.png`} />
        </div>
        <div
          style={{
            width: '500px',
            display: 'flex',
            flexDirection: 'column',
            gap: 10,
            left: '79px',
            top: '230px',
            bottom: '156px',
            position: 'absolute',
            color: 'white',
            fontSize: '50px',
            fontFamily: 'Sora',
            lineHeight: '120%',
            fontWeight: 700,
            textAlign: 'left',
          }}
        >
          {decodeURIComponent(title)}
          <div
            style={{
              color: '#A1A1AA',
              fontSize: '28px',
              fontWeight: 500,
            }}
          >
            {decodeURIComponent(description)}
          </div>
        </div>
        <br />
      </div>
    ),
    // ImageResponse options
    {
      // For convenience, we can re-use the exported opengraph-image
      // size config to also set the ImageResponse's width and height.
      ...size,
      fonts: [
        {
          name: 'Sora',
          data: await font,
          style: 'normal',
          weight: 600,
        },
      ],
    },
  )
}
