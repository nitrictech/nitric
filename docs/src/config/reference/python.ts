import { FaPython } from 'react-icons/fa'
import { NavGroup } from '../types'

export const PyReference: NavGroup = {
  title: 'Python',
  icon: FaPython,
  items: [
    {
      title: 'Overview',
      href: '/reference/python',
      breadcrumbRoot: true,
    },
    {
      title: 'APIs',
      items: [
        {
          title: 'api()',
          href: '/reference/python/api/api',
        },
        {
          title: 'api.get()',
          href: '/reference/python/api/api-get',
        },
        {
          title: 'api.post()',
          href: '/reference/python/api/api-post',
        },
        {
          title: 'api.put()',
          href: '/reference/python/api/api-put',
        },
        {
          title: 'api.delete()',
          href: '/reference/python/api/api-delete',
        },
        {
          title: 'api.patch()',
          href: '/reference/python/api/api-patch',
        },
        {
          title: 'api.methods()',
          href: '/reference/python/api/api-methods',
        },
        {
          title: 'api.all()',
          href: '/reference/python/api/api-all',
        },
      ],
    },
    {
      title: 'Batch',
      items: [
        {
          title: 'job()',
          href: '/reference/python/batch/job',
        },
        {
          title: 'job.handler()',
          href: '/reference/python/batch/job-handler',
        },
        {
          title: 'job.submit()',
          href: '/reference/python/batch/job-submit',
        },
      ],
    },
    {
      title: 'Key Value Stores',
      items: [
        {
          title: 'kv()',
          href: '/reference/python/keyvalue/keyvalue',
        },
        {
          title: 'kv.get()',
          href: '/reference/python/keyvalue/keyvalue-get',
        },
        {
          title: 'kv.set()',
          href: '/reference/python/keyvalue/keyvalue-set',
        },
        {
          title: 'kv.delete()',
          href: '/reference/python/keyvalue/keyvalue-delete',
        },
        {
          title: 'kv.keys()',
          href: '/reference/python/keyvalue/keyvalue-keys',
        },
      ],
    },
    {
      title: 'Topics',
      items: [
        {
          title: 'topic()',
          href: '/reference/python/topic/topic',
        },
        {
          title: 'topic.publish()',
          href: '/reference/python/topic/topic-publish',
        },
        {
          title: 'topic.subscribe()',
          href: '/reference/python/topic/topic-subscribe',
        },
      ],
    },
    {
      title: 'Queues',
      items: [
        {
          title: 'queue()',
          href: '/reference/python/queues/queue',
        },
        {
          title: 'queue.enqueue()',
          href: '/reference/python/queues/queue-enqueue',
        },
        {
          title: 'queue.dequeue()',
          href: '/reference/python/queues/queue-dequeue',
        },
      ],
    },
    {
      title: 'Secrets',
      items: [
        {
          title: 'secret()',
          href: '/reference/python/secrets/secret',
        },
        {
          title: 'secret.put()',
          href: '/reference/python/secrets/secret-put',
        },
        {
          title: 'secret.version()',
          href: '/reference/python/secrets/secret-version',
        },
        {
          title: 'secret.latest()',
          href: '/reference/python/secrets/secret-latest',
        },
        {
          title: 'secret.version.access()',
          href: '/reference/python/secrets/secret-version-access',
        },
      ],
    },
    {
      title: 'Storage',
      items: [
        {
          title: 'bucket()',
          href: '/reference/python/storage/bucket',
        },
        {
          title: 'bucket.on()',
          href: '/reference/python/storage/bucket-on',
        },
        {
          title: 'bucket.file()',
          href: '/reference/python/storage/bucket-file',
        },
        {
          title: 'bucket.files()',
          href: '/reference/python/storage/bucket-files',
        },
        {
          title: 'file.read()',
          href: '/reference/python/storage/bucket-file-read',
        },
        {
          title: 'file.write()',
          href: '/reference/python/storage/bucket-file-write',
        },
        {
          title: 'file.delete()',
          href: '/reference/python/storage/bucket-file-delete',
        },
        {
          title: 'file.download_url()',
          href: '/reference/python/storage/bucket-file-downloadurl',
        },
        {
          title: 'file.upload_url()',
          href: '/reference/python/storage/bucket-file-uploadurl',
        },
      ],
    },
    {
      title: 'SQL',
      items: [
        {
          title: 'sql()',
          href: '/reference/python/sql/sql',
        },
        {
          title: 'sql.connection_string()',
          href: '/reference/python/sql/sql-connection-string',
        },
      ],
    },
    {
      title: 'Schedules',
      items: [
        {
          title: 'schedule()',
          href: '/reference/python/schedule/schedule',
        },
      ],
    },
    {
      title: 'Websockets',
      items: [
        {
          title: 'websocket()',
          href: '/reference/python/websocket/websocket',
        },
        {
          title: 'websocket.on()',
          href: '/reference/python/websocket/websocket-on',
        },
        {
          title: 'websocket.send()',
          href: '/reference/python/websocket/websocket-send',
        },
      ],
    },
  ],
}
