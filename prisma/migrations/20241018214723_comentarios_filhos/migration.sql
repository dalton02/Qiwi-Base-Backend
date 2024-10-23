/*
  Warnings:

  - A unique constraint covering the columns `[tipo,titulo]` on the table `postagem` will be added. If there are existing duplicate values, this will fail.

*/
-- DropIndex
DROP INDEX "postagem_titulo_key";

-- AlterTable
ALTER TABLE "comentario" ADD COLUMN     "parenteId" INTEGER;

-- CreateIndex
CREATE UNIQUE INDEX "postagem_tipo_titulo_key" ON "postagem"("tipo", "titulo");

-- AddForeignKey
ALTER TABLE "comentario" ADD CONSTRAINT "comentario_parenteId_fkey" FOREIGN KEY ("parenteId") REFERENCES "comentario"("id") ON DELETE SET NULL ON UPDATE CASCADE;
