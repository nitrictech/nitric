import { FaPython } from 'react-icons/fa'
import { NavGroup } from '../types'
import { SiGo } from 'react-icons/si'

export const GoReference: NavGroup = {
  title: 'Go',
  icon: SiGo,
  items: [
    {
      title: 'Overview',
      href: '/reference/go',
      breadcrumbRoot: true,
    },
    {
      title: 'Resources',
      items: [
        {
          title: 'NewApi()',
          href: '/reference/go/api/api',
        },
        {
          title: 'NewJob()',
          href: '/reference/go/batch/job',
        },
        {
          title: 'NewKv()',
          href: '/reference/go/keyvalue/keyvalue',
        },
        {
          title: 'NewQueue()',
          href: '/reference/go/queues/queue',
        },
        {
          title: 'NewSecret()',
          href: '/reference/go/secrets/secret',
        },
        {
          title: 'NewSqlDatabase()',
          href: '/reference/go/sql/sql',
        },
        {
          title: 'NewBucket()',
          href: '/reference/go/storage/bucket',
        },
        {
          title: 'NewTopic()',
          href: '/reference/go/topic/topic',
        },
        {
          title: 'NewSchedule()',
          href: '/reference/go/schedule/schedule',
        },
        {
          title: 'NewWebsocket()',
          href: '/reference/go/websocket/websocket',
        },
      ],
    },
    {
      title: 'APIs',
      items: [
        {
          title: 'Api.Get()',
          href: '/reference/go/api/api-get',
        },
        {
          title: 'Api.Post()',
          href: '/reference/go/api/api-post',
        },
        {
          title: 'Api.Put()',
          href: '/reference/go/api/api-put',
        },
        {
          title: 'Api.Delete()',
          href: '/reference/go/api/api-delete',
        },
        {
          title: 'Api.Patch()',
          href: '/reference/go/api/api-patch',
        },
        {
          title: 'Api.NewRoute()',
          href: '/reference/go/api/api-route',
        },
        {
          title: 'Api.Route.All()',
          href: '/reference/go/api/api-route-all',
        },
        {
          title: 'Api.Route.Get()',
          href: '/reference/go/api/api-route-get',
        },
        {
          title: 'Api.Route.Post()',
          href: '/reference/go/api/api-route-post',
        },
        {
          title: 'Api.Route.Put()',
          href: '/reference/go/api/api-route-put',
        },
        {
          title: 'Api.Route.Delete()',
          href: '/reference/go/api/api-route-delete',
        },
        {
          title: 'Api.Route.Patch()',
          href: '/reference/go/api/api-route-patch',
        },
      ],
    },
    {
      title: 'Batch',
      items: [
        {
          title: 'Job.Handler()',
          href: '/reference/go/batch/job-handler',
        },
        {
          title: 'Job.Submit()',
          href: '/reference/go/batch/job-submit',
        },
      ],
    },
    {
      title: 'Key Value Stores',
      items: [
        {
          title: 'Kv.Get()',
          href: '/reference/go/keyvalue/keyvalue-get',
        },
        {
          title: 'Kv.Set()',
          href: '/reference/go/keyvalue/keyvalue-set',
        },
        {
          title: 'Kv.Delete()',
          href: '/reference/go/keyvalue/keyvalue-delete',
        },
        {
          title: 'Kv.Keys()',
          href: '/reference/go/keyvalue/keyvalue-keys',
        },
      ],
    },
    {
      title: 'Topics',
      items: [
        {
          title: 'Topic.Publish()',
          href: '/reference/go/topic/topic-publish',
        },
        {
          title: 'Topic.Subscribe()',
          href: '/reference/go/topic/topic-subscribe',
        },
      ],
    },
    {
      title: 'Queues',
      items: [
        {
          title: 'Queue.Enqueue()',
          href: '/reference/go/queues/queue-enqueue',
        },
        {
          title: 'Queue.Dequeue()',
          href: '/reference/go/queues/queue-dequeue',
        },
      ],
    },
    {
      title: 'Secrets',
      items: [
        {
          title: 'Secret.Put()',
          href: '/reference/go/secrets/secret-put',
        },
        {
          title: 'Secret.AccessVersion()',
          href: '/reference/go/secrets/secret-access-version',
        },
        {
          title: 'Secret.Access()',
          href: '/reference/go/secrets/secret-access',
        },
      ],
    },
    {
      title: 'Storage',
      items: [
        {
          title: 'Bucket.On()',
          href: '/reference/go/storage/bucket-on',
        },
        {
          title: 'Bucket.ListFiles()',
          href: '/reference/go/storage/bucket-listfiles',
        },
        {
          title: 'Bucket.Read()',
          href: '/reference/go/storage/bucket-read',
        },
        {
          title: 'Bucket.Write()',
          href: '/reference/go/storage/bucket-write',
        },
        {
          title: 'Bucket.Delete()',
          href: '/reference/go/storage/bucket-delete',
        },
        {
          title: 'Bucket.DownloadUrl()',
          href: '/reference/go/storage/bucket-downloadurl',
        },
        {
          title: 'Bucket.UploadUrl()',
          href: '/reference/go/storage/bucket-uploadurl',
        },
      ],
    },
    {
      title: 'SQL',
      items: [
        {
          title: 'SqlDatabase.ConnectionString()',
          href: '/reference/go/sql/sql-connection-string',
        },
      ],
    },
    {
      title: 'Schedules',
      items: [
        {
          title: 'Schedule.Every()',
          href: '/reference/go/schedule/schedule-every',
        },
        {
          title: 'Schedule.Cron()',
          href: '/reference/go/schedule/schedule-cron',
        },
      ],
    },
    {
      title: 'Websockets',
      items: [
        {
          title: 'Websocket.On()',
          href: '/reference/go/websocket/websocket-on',
        },
        {
          title: 'Websocket.Send()',
          href: '/reference/go/websocket/websocket-send',
        },
        {
          title: 'Websocket.Close()',
          href: '/reference/go/websocket/websocket-close',
        },
      ],
    },
  ],
}
