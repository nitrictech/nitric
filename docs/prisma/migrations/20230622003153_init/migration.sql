-- CreateTable
CREATE TABLE "Feedback" (
    "id" SERIAL NOT NULL,
    "url" TEXT NOT NULL,
    "answer" TEXT NOT NULL,
    "label" TEXT NOT NULL,
    "ua" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "Feedback_pkey" PRIMARY KEY ("id")
);
