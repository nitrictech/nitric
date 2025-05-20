import { api, sql } from '@nitric/sdk'
import { PrismaClient } from '@prisma/client'

const mainApi = api('main')
const db = sql('my-data')

let prisma
const getClient = async () => {
  // ensure we only create the client once
  if (!prisma) {
    const connectionString = await db.connectionString()

    prisma = new PrismaClient({
      datasourceUrl: connectionString,
    })
  }
  return prisma
}

// api demonstrating connecting to prisma and then doing an insert and select
mainApi.post('/users/:name', async (ctx) => {
  const { name } = ctx.req.params

  const client = await getClient()

  await client.user.create({
    data: {
      name,
      email: `${name}@example.com`,
    },
  })

  const createdUser = await client.user.findFirstOrThrow({
    where: {
      name,
    },
  })

  ctx.res.body = `Created ${name} with ID ${createdUser.id}`

  return ctx
})
