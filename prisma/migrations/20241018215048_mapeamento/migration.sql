/*
  Warnings:

  - You are about to drop the column `parenteId` on the `comentario` table. All the data in the column will be lost.

*/
-- DropForeignKey
ALTER TABLE "comentario" DROP CONSTRAINT "comentario_parenteId_fkey";

-- AlterTable
ALTER TABLE "comentario" DROP COLUMN "parenteId",
ADD COLUMN     "parente_id" INTEGER;

-- AddForeignKey
ALTER TABLE "comentario" ADD CONSTRAINT "comentario_parente_id_fkey" FOREIGN KEY ("parente_id") REFERENCES "comentario"("id") ON DELETE SET NULL ON UPDATE CASCADE;
