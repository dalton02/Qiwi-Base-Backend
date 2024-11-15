/*
  Warnings:

  - You are about to drop the column `codigo` on the `usuario` table. All the data in the column will be lost.
  - You are about to drop the column `curso` on the `usuario` table. All the data in the column will be lost.
  - A unique constraint covering the columns `[idAluno]` on the table `usuario` will be added. If there are existing duplicate values, this will fail.
  - Added the required column `perfil` to the `usuario` table without a default value. This is not possible if the table is not empty.

*/
-- DropForeignKey
ALTER TABLE "comentario" DROP CONSTRAINT "comentario_usuario_id_fkey";

-- DropForeignKey
ALTER TABLE "postagem" DROP CONSTRAINT "postagem_usuario_id_fkey";

-- AlterTable
ALTER TABLE "usuario" DROP COLUMN "codigo",
DROP COLUMN "curso",
ADD COLUMN     "idAluno" INTEGER,
ADD COLUMN     "perfil" VARCHAR(200) NOT NULL;

-- CreateTable
CREATE TABLE "aluno" (
    "id" SERIAL NOT NULL,
    "codigo" INTEGER NOT NULL,
    "curso" VARCHAR(200) NOT NULL,
    "idUsuario" INTEGER NOT NULL,

    CONSTRAINT "aluno_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "aluno_idUsuario_key" ON "aluno"("idUsuario");

-- CreateIndex
CREATE UNIQUE INDEX "usuario_idAluno_key" ON "usuario"("idAluno");

-- AddForeignKey
ALTER TABLE "aluno" ADD CONSTRAINT "aluno_idUsuario_fkey" FOREIGN KEY ("idUsuario") REFERENCES "usuario"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "postagem" ADD CONSTRAINT "postagem_usuario_id_fkey" FOREIGN KEY ("usuario_id") REFERENCES "aluno"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "comentario" ADD CONSTRAINT "comentario_usuario_id_fkey" FOREIGN KEY ("usuario_id") REFERENCES "aluno"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
