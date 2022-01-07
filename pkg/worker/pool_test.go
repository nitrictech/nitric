package worker

import (
	"fmt"
	"sync"

	"github.com/golang/mock/gomock"
	mock_sync "github.com/nitrictech/nitric/mocks/sync"
	mock_worker "github.com/nitrictech/nitric/mocks/worker"
	"github.com/nitrictech/nitric/pkg/triggers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProcessPool", func() {

	Context("GetWorkerCount", func() {

		When("calling GetWorkerCount", func() {
			ctrl := gomock.NewController(GinkgoT())
			lck := mock_sync.NewMockLocker(ctrl)

			pp := &ProcessPool{
				workerLock: lck,
				workers:    make([]Worker, 0),
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
				hw := &RouteWorker{}
				ew := &SubscriptionWorker{}
				fw := &FaasWorker{}

				pp := &ProcessPool{
					maxWorkers: 3,
					workerLock: &sync.Mutex{},
					workers:    []Worker{hw, ew, fw},
				}

				wrkrs := pp.getHttpWorkers()

				It("should return all http capable workers", func() {
					Expect(wrkrs).To(HaveLen(2))
				})

				It("should prioritise route workers", func() {
					Expect(wrkrs[0]).To(Equal(hw))
				})

				It("should return other http capable workers", func() {
					Expect(wrkrs[1]).To(Equal(fw))
				})
			})
		})

		Context("getEventWorkers", func() {
			When("pool contains mix of event & http handlers", func() {
				hw := &RouteWorker{}
				ew := &SubscriptionWorker{}
				fw := &FaasWorker{}

				pp := &ProcessPool{
					maxWorkers: 3,
					workerLock: &sync.Mutex{},
					workers:    []Worker{hw, ew, fw},
				}

				wrkrs := pp.getEventWorkers()

				It("should return all event capable workers", func() {
					Expect(wrkrs).To(HaveLen(2))
				})

				It("should prioritise specialized workers", func() {
					Expect(wrkrs[0]).To(Equal(ew))
				})

				It("should return other event capable workers", func() {
					Expect(wrkrs[1]).To(Equal(fw))
				})
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
				pp := &ProcessPool{minWorkers: 1, workers: make([]Worker, 0), workerLock: &sync.Mutex{}}

				It("should return an error", func() {
					err := pp.WaitForMinimumWorkers(0)

					Expect(err).Should(HaveOccurred())
				})
			})

			When("minimum worker count is eventually met", func() {
				pp := &ProcessPool{minWorkers: 1, workers: make([]Worker, 0), workerLock: &sync.Mutex{}}

				It("should block the current thread until it is", func() {
					wg := sync.WaitGroup{}
					wg.Add(1)
					var err error

					go func() {
						defer wg.Done()
						err = pp.WaitForMinimumWorkers(100)
					}()

					pp.AddWorker(&FaasWorker{})

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
					badWrkr := mock_worker.NewMockGrpcWorker(ctrl)
					pp := &ProcessPool{minWorkers: 0, workers: []Worker{badWrkr}, workerLock: &sync.Mutex{}}

					It("should return an error", func() {
						By("testing the worker with the trigger")
						badWrkr.EXPECT().HandlesHttpRequest(gomock.Any()).Return(false).Times(1)

						By("returning a nil worker")
						wrkr, err := pp.GetWorker(&GetWorkerOptions{Http: &triggers.HttpRequest{}})
						Expect(wrkr).To(BeNil())

						By("return an error")
						Expect(err).Should(HaveOccurred())
					})
				})

				When("compatible workers are available", func() {
					ctrl := gomock.NewController(GinkgoT())
					hw := mock_worker.NewMockGrpcWorker(ctrl)
					pp := &ProcessPool{minWorkers: 0, workers: []Worker{hw}, workerLock: &sync.Mutex{}}
					tr := &triggers.HttpRequest{}

					It("should return a compatible worker", func() {
						By("Querying testing the worker with the trigger")
						hw.EXPECT().HandlesHttpRequest(tr).Return(true).Times(1)

						By("returning a nil worker")
						wrkr, err := pp.GetWorker(&GetWorkerOptions{Http: tr})
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
						badWrkr := mock_worker.NewMockGrpcWorker(ctrl)
						pp := &ProcessPool{minWorkers: 0, workers: []Worker{badWrkr}, workerLock: &sync.Mutex{}}

						It("should return an error", func() {
							By("testing the worker with the trigger")
							badWrkr.EXPECT().HandlesEvent(gomock.Any()).Return(false).Times(1)

							By("returning a nil worker")
							wrkr, err := pp.GetWorker(&GetWorkerOptions{Event: &triggers.Event{}})
							Expect(wrkr).To(BeNil())

							By("return an error")
							Expect(err).Should(HaveOccurred())
						})
					})

					When("compatible workers are available", func() {
						ctrl := gomock.NewController(GinkgoT())
						hw := mock_worker.NewMockGrpcWorker(ctrl)
						pp := &ProcessPool{minWorkers: 0, workers: []Worker{hw}, workerLock: &sync.Mutex{}}
						tr := &triggers.Event{}

						It("should return a compatible worker", func() {
							By("Querying testing the worker with the trigger")
							hw.EXPECT().HandlesEvent(tr).Return(true).Times(1)

							By("returning a nil worker")
							wrkr, err := pp.GetWorker(&GetWorkerOptions{Event: tr})
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
				wkr := mock_worker.NewMockGrpcWorker(ctrl)

				pp := &ProcessPool{
					workerLock: lck,
					workers:    []Worker{wkr},
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
				wkr := mock_worker.NewMockGrpcWorker(ctrl)

				pp := &ProcessPool{
					workerLock: lck,
					workers:    []Worker{&FaasWorker{}},
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
				wkr := mock_worker.NewMockGrpcWorker(ctrl)

				pp := &ProcessPool{
					workerLock: lck,
					workers:    []Worker{},
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
				wkr := mock_worker.NewMockGrpcWorker(ctrl)
				wkr2 := mock_worker.NewMockGrpcWorker(ctrl)

				pp := &ProcessPool{
					workerLock: lck,
					workers:    []Worker{wkr},
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

		})
	})
})
