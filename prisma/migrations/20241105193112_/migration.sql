/*
  Warnings:

  - You are about to drop the column `idUsuario` on the `aluno` table. All the data in the column will be lost.
  - You are about to drop the column `idAluno` on the `usuario` table. All the data in the column will be lost.
  - A unique constraint covering the columns `[id_usuario]` on the table `aluno` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[id_aluno]` on the table `usuario` will be added. If there are existing duplicate values, this will fail.
  - Added the required column `id_usuario` to the `aluno` table without a default value. This is not possible if the table is not empty.

*/
-- DropForeignKey
ALTER TABLE "aluno" DROP CONSTRAINT "aluno_idUsuario_fkey";

-- DropIndex
DROP INDEX "aluno_idUsuario_key";

-- DropIndex
DROP INDEX "usuario_idAluno_key";

-- AlterTable
ALTER TABLE "aluno" DROP COLUMN "idUsuario",
ADD COLUMN     "id_usuario" INTEGER NOT NULL;

-- AlterTable
ALTER TABLE "usuario" DROP COLUMN "idAluno",
ADD COLUMN     "id_aluno" INTEGER;

-- CreateIndex
CREATE UNIQUE INDEX "aluno_id_usuario_key" ON "aluno"("id_usuario");

-- CreateIndex
CREATE UNIQUE INDEX "usuario_id_aluno_key" ON "usuario"("id_aluno");

-- AddForeignKey
ALTER TABLE "aluno" ADD CONSTRAINT "aluno_id_usuario_fkey" FOREIGN KEY ("id_usuario") REFERENCES "usuario"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
