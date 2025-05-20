import { FaNodeJs } from 'react-icons/fa'
import { NavGroup } from '../types'

export const NodeReference: NavGroup = {
  title: 'Node.js',
  icon: FaNodeJs,
  items: [
    {
      title: 'Overview',
      href: '/reference/nodejs',
      breadcrumbRoot: true,
    },
    {
      title: 'APIs',
      items: [
        {
          title: 'api()',
          href: '/reference/nodejs/api/api',
        },
        {
          title: 'api.get()',
          href: '/reference/nodejs/api/api-get',
        },
        {
          title: 'api.post()',
          href: '/reference/nodejs/api/api-post',
        },
        {
          title: 'api.put()',
          href: '/reference/nodejs/api/api-put',
        },
        {
          title: 'api.delete()',
          href: '/reference/nodejs/api/api-delete',
        },
        {
          title: 'api.patch()',
          href: '/reference/nodejs/api/api-patch',
        },
        {
          title: 'api.route()',
          href: '/reference/nodejs/api/api-route',
        },
        {
          title: 'api.route.all()',
          href: '/reference/nodejs/api/api-route-all',
        },
        {
          title: 'api.route.get()',
          href: '/reference/nodejs/api/api-route-get',
        },
        {
          title: 'api.route.post()',
          href: '/reference/nodejs/api/api-route-post',
        },
        {
          title: 'api.route.put()',
          href: '/reference/nodejs/api/api-route-put',
        },
        {
          title: 'api.route.delete()',
          href: '/reference/nodejs/api/api-route-delete',
        },
        {
          title: 'api.route.patch()',
          href: '/reference/nodejs/api/api-route-patch',
        },
      ],
    },
    {
      title: 'Batch',
      items: [
        {
          title: 'job()',
          href: '/reference/nodejs/batch/job',
        },
        {
          title: 'job.handler()',
          href: '/reference/nodejs/batch/job-handler',
        },
        {
          title: 'job.submit()',
          href: '/reference/nodejs/batch/job-submit',
        },
      ],
    },
    {
      title: 'HTTP',
      items: [
        {
          title: 'http()',
          href: '/reference/nodejs/http/http',
        },
      ],
    },
    {
      title: 'Key Value Stores',
      items: [
        {
          title: 'kv()',
          href: '/reference/nodejs/keyvalue/keyvalue',
        },
        {
          title: 'kv.get()',
          href: '/reference/nodejs/keyvalue/keyvalue-get',
        },
        {
          title: 'kv.set()',
          href: '/reference/nodejs/keyvalue/keyvalue-set',
        },
        {
          title: 'kv.delete()',
          href: '/reference/nodejs/keyvalue/keyvalue-delete',
        },
        {
          title: 'kv.keys()',
          href: '/reference/nodejs/keyvalue/keyvalue-keys',
        },
      ],
    },
    {
      title: 'Topics',
      items: [
        {
          title: 'topic()',
          href: '/reference/nodejs/topic/topic',
        },
        {
          title: 'topic.publish()',
          href: '/reference/nodejs/topic/topic-publish',
        },
        {
          title: 'topic.subscribe()',
          href: '/reference/nodejs/topic/topic-subscribe',
        },
      ],
    },
    {
      title: 'Queues',
      items: [
        {
          title: 'queue()',
          href: '/reference/nodejs/queues/queue',
        },
        {
          title: 'queue.enqueue()',
          href: '/reference/nodejs/queues/queue-enqueue',
        },
        {
          title: 'queue.dequeue()',
          href: '/reference/nodejs/queues/queue-dequeue',
        },
      ],
    },
    {
      title: 'Secrets',
      items: [
        {
          title: 'secret()',
          href: '/reference/nodejs/secrets/secret',
        },
        {
          title: 'secret.put()',
          href: '/reference/nodejs/secrets/secret-put',
        },
        {
          title: 'secret.version()',
          href: '/reference/nodejs/secrets/secret-version',
        },
        {
          title: 'secret.latest()',
          href: '/reference/nodejs/secrets/secret-latest',
        },
        {
          title: 'secret.version.access()',
          href: '/reference/nodejs/secrets/secret-version-access',
        },
      ],
    },
    {
      title: 'Storage',
      items: [
        {
          title: 'bucket()',
          href: '/reference/nodejs/storage/bucket',
        },
        {
          title: 'bucket.on()',
          href: '/reference/nodejs/storage/bucket-on',
        },
        {
          title: 'bucket.file()',
          href: '/reference/nodejs/storage/bucket-file',
        },
        {
          title: 'bucket.files()',
          href: '/reference/nodejs/storage/bucket-files',
        },
        {
          title: 'file.exists()',
          href: '/reference/nodejs/storage/bucket-file-exists',
        },
        {
          title: 'file.read()',
          href: '/reference/nodejs/storage/bucket-file-read',
        },
        {
          title: 'file.write()',
          href: '/reference/nodejs/storage/bucket-file-write',
        },
        {
          title: 'file.delete()',
          href: '/reference/nodejs/storage/bucket-file-delete',
        },
        {
          title: 'file.getDownloadUrl()',
          href: '/reference/nodejs/storage/bucket-file-downloadurl',
        },
        {
          title: 'file.getUploadUrl()',
          href: '/reference/nodejs/storage/bucket-file-uploadurl',
        },
      ],
    },
    {
      title: 'SQL',
      items: [
        {
          title: 'sql()',
          href: '/reference/nodejs/sql/sql',
        },
        {
          title: 'sql.connectionString()',
          href: '/reference/nodejs/sql/sql-connection-string',
        },
      ],
    },
    {
      title: 'Schedules',
      items: [
        {
          title: 'schedule()',
          href: '/reference/nodejs/schedule/schedule',
        },
        {
          title: 'schedule.every()',
          href: '/reference/nodejs/schedule/schedule-every',
        },
        {
          title: 'schedule.cron()',
          href: '/reference/nodejs/schedule/schedule-cron',
        },
      ],
    },
    {
      title: 'Websockets',
      items: [
        {
          title: 'websocket()',
          href: '/reference/nodejs/websocket/websocket',
        },
        {
          title: 'websocket.on()',
          href: '/reference/nodejs/websocket/websocket-on',
        },
        {
          title: 'websocket.send()',
          href: '/reference/nodejs/websocket/websocket-send',
        },
        {
          title: 'websocket.close()',
          href: '/reference/nodejs/websocket/websocket-close',
        },
      ],
    },
  ],
}
