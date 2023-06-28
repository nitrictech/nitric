// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pool

import (
	"fmt"
	"sync"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_sync "github.com/nitrictech/nitric/core/mocks/sync"
	mock_worker "github.com/nitrictech/nitric/core/mocks/worker"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/worker"
)

var _ = Describe("ProcessPool", func() {
	Context("NewProcessPool", func() {
		When("calling NewProcessPool", func() {
			pp := NewProcessPool(&ProcessPoolOptions{
				MinWorkers: 1,
				MaxWorkers: 3,
			})

			It("should have the correct options set", func() {
				Expect(pp.GetMinWorkers()).To(Equal(1))
				Expect(pp.GetMaxWorkers()).To(Equal(3))
				Expect(pp.GetWorkerCount()).To(Equal(0))
			})

			It("should not return an error", func() {
				Expect(pp).To(Not(BeNil()))
			})
		})

		When("calling NewProcessPool with max workers less than 1", func() {
			pp := NewProcessPool(&ProcessPoolOptions{
				MinWorkers: 0,
				MaxWorkers: 0,
			})

			It("should have 1 max worker", func() {
				Expect(pp.GetMaxWorkers()).To(Equal(1))
			})

			It("should not return an error", func() {
				Expect(pp).To(Not(BeNil()))
			})
		})
	})
	Context("GetWorkerCount", func() {
		When("calling GetWorkerCount", func() {
			ctrl := gomock.NewController(GinkgoT())
			lck := mock_sync.NewMockLocker(ctrl)

			pp := &ProcessPool{
				workerLock: lck,
				workers:    make([]worker.Worker, 0),
				poolErr:    make(chan error),
			}

			It("should thread safely return the number of workers", func() {
				By("Locking the worker lock")
				lck.EXPECT().Lock().Times(1)

				By("Unlocking the worker lock")
				lck.EXPECT().Unlock().Times(1)

				By("Returning the number of workers")
				Expect(pp.GetWorkerCount()).To(Equal(0))

				ctrl.Finish()
			})
		})

		Context("getHttpWorkers", func() {
			When("pool contains mix of event & http handlers", func() {
				hw := &worker.RouteWorker{}
				ew := &worker.SubscriptionWorker{}
				pw := &worker.HttpWorker{}
				fw := &worker.FaasWorker{}

				pp := &ProcessPool{
					maxWorkers: 3,
					workerLock: &sync.Mutex{},
					workers:    []worker.Worker{pw, hw, ew, fw},
				}

				wrkrs := pp.getHttpWorkers()

				It("should return all http capable workers", func() {
					Expect(wrkrs).To(HaveLen(3))
				})

				It("should prioritise route workers and http workers", func() {
					Expect(wrkrs[0]).To(Equal(hw))
					Expect(wrkrs[1]).To(Equal(pw))
				})

				It("should return other http capable workers", func() {
					Expect(wrkrs[2]).To(Equal(fw))
				})
			})
		})

		Context("getEventWorkers", func() {
			When("pool contains mix of event & http handlers", func() {
				hw := &worker.RouteWorker{}
				ew := &worker.SubscriptionWorker{}
				sw := &worker.ScheduleWorker{}
				fw := &worker.FaasWorker{}

				pp := &ProcessPool{
					maxWorkers: 4,
					workerLock: &sync.Mutex{},
					workers:    []worker.Worker{hw, sw, ew, fw},
				}

				wrkrs := pp.getEventWorkers()

				It("should return all event capable workers", func() {
					Expect(wrkrs).To(HaveLen(3))
				})

				It("should prioritise specialized workers", func() {
					Expect(wrkrs[0]).To(Equal(ew))
					Expect(wrkrs[1]).To(Equal(sw))
				})

				It("should return other event capable workers", func() {
					Expect(wrkrs[2]).To(Equal(fw))
				})
			})
		})

		Context("getNotificationWorkers", func() {
			When("pool contains mix of ", func() {
				nw := &worker.BucketNotificationWorker{}
				sw := &worker.ScheduleWorker{}
				fw := &worker.FaasWorker{}

				pp := &ProcessPool{
					maxWorkers: 3,
					workerLock: &sync.Mutex{},
					workers:    []worker.Worker{fw, nw, sw},
				}

				wrkrs := pp.getNotificationWorkers()

				It("should return all notification capable workers", func() {
					Expect(wrkrs).To(HaveLen(1))
				})

				It("should prioritise specialized workers", func() {
					Expect(wrkrs[0]).To(Equal(nw))
				})
			})
		})
	})

	Context("GetWorkers", func() {
		When("calling without filters", func() {
			nw := &worker.BucketNotificationWorker{}
			sw := &worker.ScheduleWorker{}
			fw := &worker.FaasWorker{}

			pp := &ProcessPool{
				maxWorkers: 3,
				workerLock: &sync.Mutex{},
				workers:    []worker.Worker{fw, nw, sw},
			}

			wrkrs := pp.GetWorkers(nil)

			It("should return all workers", func() {
				Expect(wrkrs).To(HaveLen(3))
			})
		})

		When("calling with filters", func() {
			rw := &worker.RouteWorker{}
			nw := &worker.BucketNotificationWorker{}
			tw := &worker.SubscriptionWorker{}
			sw := &worker.ScheduleWorker{}
			fw := &worker.FaasWorker{}

			pp := &ProcessPool{
				maxWorkers: 5,
				workerLock: &sync.Mutex{},
				workers:    []worker.Worker{rw, tw, fw, nw, sw},
			}

			wrkrs := pp.GetWorkers(&GetWorkerOptions{
				Filter: func(w worker.Worker) bool {
					_, ok := w.(*worker.RouteWorker)
					return ok
				},
			})

			It("should return only route workers", func() {
				Expect(wrkrs).To(HaveLen(1))
			})

			It("should return", func() {
				Expect(wrkrs[0]).To(Equal(rw))
			})
		})

		Context("GetMinWorkers", func() {
			When("calling getMinWorkers", func() {
				pp := &ProcessPool{minWorkers: 12}

				It("should return the value of the minWorkers field", func() {
					Expect(pp.GetMinWorkers()).To(Equal(pp.minWorkers))
				})
			})
		})

		Context("GetMaxWorkers", func() {
			When("calling getMaxWorkers", func() {
				pp := &ProcessPool{maxWorkers: 20}
				It("should return the value of the maxWorkers field", func() {
					Expect(pp.GetMaxWorkers()).To(Equal(pp.maxWorkers))
				})
			})
		})

		Context("Monitor", func() {
			When("calling monitor", func() {
				poolChan := make(chan error)
				pp := &ProcessPool{
					poolErr: poolChan,
				}
				mockErr := fmt.Errorf("mock error")

				It("should block the current thread until poolErr is populated", func() {
					By("Blocking until pool error is populated")
					blockingErr := make(chan error)
					go func(errChan chan error) {
						errChan <- pp.Monitor()
					}(blockingErr)

					// populate poolChan
					poolChan <- mockErr

					// Blocking here (should immediately unblock as blockingErr == <- poolChan)
					err := <-blockingErr

					By("Capturing the original pool error")
					Expect(err).To(Equal(mockErr))
				})
			})
		})

		Context("WaitForMinimumWorkers", func() {
			When("minimum worker count is not met before timeout", func() {
				pp := &ProcessPool{minWorkers: 1, workers: make([]worker.Worker, 0), workerLock: &sync.Mutex{}}

				It("should return an error", func() {
					err := pp.WaitForMinimumWorkers(0)

					Expect(err).Should(HaveOccurred())
				})
			})

			When("minimum worker count is eventually met", func() {
				pp := &ProcessPool{minWorkers: 1, workers: make([]worker.Worker, 0), workerLock: &sync.Mutex{}}

				It("should block the current thread until it is", func() {
					wg := sync.WaitGroup{}
					wg.Add(1)
					var err error

					go func() {
						defer wg.Done()
						err = pp.WaitForMinimumWorkers(100)
					}()

					err = pp.AddWorker(&worker.FaasWorker{})
					Expect(err).To(BeNil())

					By("waiting for the worker")
					wg.Wait()

					By("not returning an error")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})

		Context("GetWorker", func() {
			Context("Getting a worker for a Http trigger", func() {
				When("no compatible workers are available", func() {
					ctrl := gomock.NewController(GinkgoT())
					badWrkr := mock_worker.NewMockWorker(ctrl)
					pp := &ProcessPool{minWorkers: 0, workers: []worker.Worker{badWrkr}, workerLock: &sync.Mutex{}}

					It("should return an error", func() {
						By("testing the worker with the trigger")
						badWrkr.EXPECT().HandlesTrigger(gomock.Any()).Return(false).Times(1)

						By("returning a nil worker")
						wrkr, err := pp.GetWorker(&GetWorkerOptions{Trigger: &v1.TriggerRequest{}})
						Expect(wrkr).To(BeNil())

						By("return an error")
						Expect(err).Should(HaveOccurred())
					})
				})

				When("compatible workers are available", func() {
					ctrl := gomock.NewController(GinkgoT())
					hw := mock_worker.NewMockWorker(ctrl)
					pp := &ProcessPool{minWorkers: 0, workers: []worker.Worker{hw}, workerLock: &sync.Mutex{}}
					tr := &v1.TriggerRequest{}

					It("should return a compatible worker", func() {
						By("Querying testing the worker with the trigger")
						hw.EXPECT().HandlesTrigger(tr).Return(true).Times(1)

						By("not returning a nil worker")
						wrkr, err := pp.GetWorker(&GetWorkerOptions{Trigger: tr})
						Expect(wrkr).To(Equal(hw))

						By("not returning an error")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Context("Getting a worker for an Event trigger", func() {
				When("no compatible event workers are available", func() {
					When("no compatible workers are available", func() {
						ctrl := gomock.NewController(GinkgoT())
						badWrkr := mock_worker.NewMockWorker(ctrl)
						pp := &ProcessPool{minWorkers: 0, workers: []worker.Worker{badWrkr}, workerLock: &sync.Mutex{}}

						It("should return an error", func() {
							By("testing the worker with the trigger")
							badWrkr.EXPECT().HandlesTrigger(gomock.Any()).Return(false).Times(1)

							By("returning a nil worker")
							wrkr, err := pp.GetWorker(&GetWorkerOptions{Trigger: &v1.TriggerRequest{}})
							Expect(wrkr).To(BeNil())

							By("return an error")
							Expect(err).Should(HaveOccurred())
						})
					})

					When("compatible workers are available", func() {
						ctrl := gomock.NewController(GinkgoT())
						hw := mock_worker.NewMockWorker(ctrl)
						pp := &ProcessPool{minWorkers: 0, workers: []worker.Worker{hw}, workerLock: &sync.Mutex{}}
						tr := &v1.TriggerRequest{}

						It("should return a compatible worker", func() {
							By("Querying testing the worker with the trigger")
							hw.EXPECT().HandlesTrigger(tr).Return(true).Times(1)

							By("returning a nil worker")
							wrkr, err := pp.GetWorker(&GetWorkerOptions{Trigger: tr})
							Expect(wrkr).To(Equal(hw))

							By("not returning an error")
							Expect(err).ShouldNot(HaveOccurred())
						})
					})
				})
			})

			Context("Getting a worker for a Bucket notification trigger", func() {
				When("no compatible event workers are available", func() {
					When("no compatible workers are available", func() {
						ctrl := gomock.NewController(GinkgoT())
						badWrkr := mock_worker.NewMockWorker(ctrl)
						pp := &ProcessPool{minWorkers: 0, workers: []worker.Worker{badWrkr}, workerLock: &sync.Mutex{}}

						It("should return an error", func() {
							By("testing the worker with the trigger")
							badWrkr.EXPECT().HandlesTrigger(gomock.Any()).Return(false).Times(1)

							By("returning a nil worker")
							wrkr, err := pp.GetWorker(&GetWorkerOptions{Trigger: &v1.TriggerRequest{}})
							Expect(wrkr).To(BeNil())

							By("return an error")
							Expect(err).Should(HaveOccurred())
						})
					})

					When("compatible workers are available", func() {
						ctrl := gomock.NewController(GinkgoT())
						hw := mock_worker.NewMockWorker(ctrl)
						pp := &ProcessPool{minWorkers: 0, workers: []worker.Worker{hw}, workerLock: &sync.Mutex{}}
						tr := &v1.TriggerRequest{}

						It("should return a compatible worker", func() {
							By("Querying testing the worker with the trigger")
							hw.EXPECT().HandlesTrigger(tr).Return(true).Times(1)

							By("returning a nil worker")
							wrkr, err := pp.GetWorker(&GetWorkerOptions{Trigger: tr})
							Expect(wrkr).To(Equal(hw))

							By("not returning an error")
							Expect(err).ShouldNot(HaveOccurred())
						})
					})
				})
			})
		})

		Context("RemoveWorker", func() {
			When("removing an existing worker from the pool", func() {
				ctrl := gomock.NewController(GinkgoT())
				lck := mock_sync.NewMockLocker(ctrl)
				wkr := mock_worker.NewMockWorker(ctrl)

				pp := &ProcessPool{
					workerLock: lck,
					workers:    []worker.Worker{wkr},
					minWorkers: 0,
					maxWorkers: 1,
				}

				It("should thread safely remove the worker", func() {
					By("locking the worker lock")
					lck.EXPECT().Lock().Times(1)

					By("unlocking the worker lock")
					lck.EXPECT().Unlock().Times(1)

					By("removing the worker from the pool")
					err := pp.RemoveWorker(wkr)

					By("not returning an error")
					Expect(err).To(BeNil())

					ctrl.Finish()
				})
			})

			When("removing a non-existent worker from the pool", func() {
				ctrl := gomock.NewController(GinkgoT())
				lck := mock_sync.NewMockLocker(ctrl)
				wkr := mock_worker.NewMockWorker(ctrl)

				pp := &ProcessPool{
					workerLock: lck,
					workers:    []worker.Worker{&worker.FaasWorker{}},
					minWorkers: 0,
					maxWorkers: 1,
				}

				It("should thread safely return an error", func() {
					By("locking the worker lock")
					lck.EXPECT().Lock().Times(1)

					By("unlocking the worker lock")
					lck.EXPECT().Unlock().Times(1)

					By("returning an error")
					Expect(pp.RemoveWorker(wkr)).Should(HaveOccurred())
				})
			})
		})

		Context("AddWorker", func() {
			When("Max workers have not been exceeded", func() {
				ctrl := gomock.NewController(GinkgoT())
				lck := mock_sync.NewMockLocker(ctrl)
				wkr := mock_worker.NewMockWorker(ctrl)

				pp := &ProcessPool{
					workerLock: lck,
					workers:    []worker.Worker{},
					minWorkers: 0,
					maxWorkers: 1,
				}

				It("should thread safely add the worker", func() {
					By("locking the worker lock")
					lck.EXPECT().Lock().Times(1)

					By("unlocking the worker lock")
					lck.EXPECT().Unlock().Times(1)

					By("adding the worker to the pool")
					err := pp.AddWorker(wkr)
					Expect(pp.workers).To(HaveLen(1))

					By("not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					ctrl.Finish()
				})
			})

			When("Max workers have been exceeded", func() {
				ctrl := gomock.NewController(GinkgoT())
				lck := mock_sync.NewMockLocker(ctrl)
				wkr := mock_worker.NewMockWorker(ctrl)
				wkr2 := mock_worker.NewMockWorker(ctrl)

				pp := &ProcessPool{
					workerLock: lck,
					workers:    []worker.Worker{wkr},
					minWorkers: 0,
					maxWorkers: 1,
				}

				It("should thread safely return an error", func() {
					By("locking the worker lock")
					lck.EXPECT().Lock().Times(1)

					By("unlocking the worker lock")
					lck.EXPECT().Unlock().Times(1)

					By("removing the worker from the pool")
					err := pp.RemoveWorker(wkr2)

					By("not returning an error")
					Expect(err).Should(HaveOccurred())

					ctrl.Finish()
				})
			})

			Context("adding bucket notification workers", func() {
				When("there are no conflicting bucket notifications", func() {
					ctrl := gomock.NewController(GinkgoT())
					lck := mock_sync.NewMockLocker(ctrl)

					nw1 := worker.NewBucketNotificationWorker(nil, &worker.BucketNotificationWorkerOptions{
						Notification: &v1.BucketNotificationWorker{
							Bucket: "test",
							Config: &v1.BucketNotificationConfig{
								NotificationType:         v1.BucketNotificationType_Created,
								NotificationPrefixFilter: "test",
							},
						},
					})

					nw2 := worker.NewBucketNotificationWorker(nil, &worker.BucketNotificationWorkerOptions{
						Notification: &v1.BucketNotificationWorker{
							Bucket: "test",
							Config: &v1.BucketNotificationConfig{
								NotificationType:         v1.BucketNotificationType_Created,
								NotificationPrefixFilter: "other",
							},
						},
					})

					pp := &ProcessPool{
						workerLock: lck,
						workers:    []worker.Worker{nw1},
						minWorkers: 0,
						maxWorkers: 1,
					}

					It("should thread safely add the worker", func() {
						By("locking the worker lock")
						lck.EXPECT().Lock().Times(1)

						By("unlocking the worker lock")
						lck.EXPECT().Unlock().Times(1)

						By("adding the worker to the pool")
						err := pp.AddWorker(nw2)
						Expect(pp.workers).To(HaveLen(2))

						By("not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						ctrl.Finish()
					})
				})

				When("there are conflicting bucket notifications", func() {
					ctrl := gomock.NewController(GinkgoT())
					lck := mock_sync.NewMockLocker(ctrl)

					nw1 := worker.NewBucketNotificationWorker(nil, &worker.BucketNotificationWorkerOptions{
						Notification: &v1.BucketNotificationWorker{
							Bucket: "test",
							Config: &v1.BucketNotificationConfig{
								NotificationType:         v1.BucketNotificationType_Created,
								NotificationPrefixFilter: "test",
							},
						},
					})

					nw2 := worker.NewBucketNotificationWorker(nil, &worker.BucketNotificationWorkerOptions{
						Notification: &v1.BucketNotificationWorker{
							Bucket: "test",
							Config: &v1.BucketNotificationConfig{
								NotificationType:         v1.BucketNotificationType_Created,
								NotificationPrefixFilter: "test/dog.png",
							},
						},
					})

					pp := &ProcessPool{
						workerLock: lck,
						workers:    []worker.Worker{nw1},
						minWorkers: 0,
						maxWorkers: 1,
					}

					It("should thread safely add the worker", func() {
						By("locking the worker lock")
						lck.EXPECT().Lock().Times(1)

						By("unlocking the worker lock")
						lck.EXPECT().Unlock().Times(1)

						By("adding the worker to the pool")
						err := pp.AddWorker(nw2)
						Expect(pp.workers).To(HaveLen(1))

						By("returning an error")
						Expect(err).Should(HaveOccurred())

						ctrl.Finish()
					})
				})
			})
		})
	})
})
