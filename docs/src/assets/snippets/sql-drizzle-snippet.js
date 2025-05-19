import { api, sql } from '@nitric/sdk'
import { drizzle } from 'drizzle-orm/postgres-js'
import * as schema from '../schema'
import postgres from 'postgres'
import { eq } from 'drizzle-orm'

const mainApi = api('main')
const db = sql('my-data')

let drizzleClient
const getClient = async () => {
  // ensure we only create the client once
  if (!drizzleClient) {
    const connectionString = await db.connectionString()

    const queryClient = postgres(connectionString)
    drizzleClient = drizzle(queryClient, { schema })
  }
  return drizzleClient
}

// api demonstrating connecting with drizzle and then doing an insert and select
mainApi.post('/users/:name', async (ctx) => {
  const { name } = ctx.req.params

  const client = await getClient()

  await client.insert(schema.users).values({
    name,
    email: `${name}@example.com`,
  })

  const createdUser = await client
    .select()
    .from(schema.users)
    .where(eq(schema.users.name, name))

  ctx.res.body = `Created ${name} with ID ${createdUser[0].id}`

  return ctx
})
