import sys

import numba
import json
import mysql.connector
import numpy as np
import time
import sys

def crack():

    cnx = mysql.connector.connect(user='root', password='a52a52a52', host='127.0.0.1', database='equationguesses')
    cursor = cnx.cursor()

    for i in range(0,131072):
        print("Trying R4 = ", i)

        query = ("SELECT * FROM (SELECT * FROM equationguesses.equation WHERE R4 = " + str(i) + " LIMIT 570) a ORDER BY Equationnumber")

        cursor.execute(query)

        M = np.zeros((570,720),dtype=np.uint8)

        idx = 0
        for (_, equation, _,_) in cursor:
            M[idx] = np.array(list(equation),dtype=np.uint8)

            idx += 1

        file = open(sys.argv[1], "r")
        content = file.read()
        known_keystream = np.array(list(map(int,content)))
        M[:,-1] = np.bitwise_xor(M[:,-1],known_keystream)

        xs = M[:,0:19]
        ys = M[:,190:212]
        zs = M[:,443:466]

        l_xs = M[:,19:190]
        l_ys = M[:,212:443]
        l_zs = M[:,466:719]

        new_M = np.append(xs,ys,axis=1)
        new_M = np.append(new_M, zs, axis=1)
        new_M = np.append(new_M, l_xs, axis=1)
        new_M = np.append(new_M, l_ys, axis=1)
        new_M = np.append(new_M, l_zs, axis=1)
        new_M = np.append(new_M, M[:,719:720], axis=1)

        solvedEquations = gf2elim(new_M)

        try:
            # Check if we have solved the equation for all non-linearized variables
            for idx,val in enumerate(solvedEquations[:64,:64]):
                if solvedEquations[idx,idx] != 1 or np.sum(solvedEquations[idx,:-1]) != 1:
                    raise AssertionError

            solutions = solvedEquations[:,-1]

            indexMapR1, luR1 = MakeIndexMapAndLookup(19)
            indexMapR2, luR2 = MakeIndexMapAndLookup(22)
            indexMapR3, luR3 = MakeIndexMapAndLookup(23)

            # Check linearised variables are correct
            for idx1, row in enumerate(solvedEquations[64:570,:]):
                result = 0
                for idx2, val in enumerate(row[:64]):
                    if val == 1:
                        result = result ^ solutions[idx2]

                for idx2, val in enumerate(row[64:235]):
                    if val == 1:
                        x1,x2 = luR1[idx2]
                        result = result ^ solutions[x1]*solutions[x2]

                for idx2, val in enumerate(row[235:466]):
                    if val == 1:
                        y1,y2 = luR2[idx2]
                        result = result ^ solutions[y1+19]*solutions[y2+19]

                for idx2, val in enumerate(row[466:719]):
                    if val == 1:
                        z1,z2 = luR3[idx2]
                        result = result ^ solutions[z1+19+22]*solutions[z2+19+22]

                assert result == solutions[idx1+64]

            # Assert forced bits are correct
            assert solutions[15] == 1
            assert solutions[19+16] == 1
            assert solutions[19+22+18] == 1
        except AssertionError:
            continue
        print("Correct R4 is: ", i)
        print("Internal state of r1 at first frame after initialization: ", solutions[:19])
        print("Internal state of r2 at first frame after initialization: ", solutions[19:19+22])
        print("Internal state of r3 at first frame after initialization: ", solutions[19+22:19+22+23])

        outFile = open("Initial states.txt", "a")
        outFile.truncate(0)
        outFile.write(np.array2string(solutions[:19],separator="")[1:-1])
        outFile.write("\n")
        outFile.write(np.array2string(solutions[19:19+22], separator="")[1:-1])
        outFile.write("\n")
        outFile.write(np.array2string(solutions[19+22:19+22+23], separator="")[1:-1])
        outFile.write("\n")
        outFile.write(str(i))
        outFile.close()

        break

    cursor.close()
    cnx.close()

# maps linearized values to a map
def MakeIndexMapAndLookup(len):
    reverseLookUp = np.zeros(int((len*(len-1))/2), dtype='i,i')
    indexLookup = 0
    index = len

    indexMap = np.zeros((len,len))

    for i in range(len):
        for j in range (i+1,len):
            indexMap[i][j] = index
            reverseLookUp[indexLookup] = (i,j)
            index += 1
            indexLookup += 1

    return indexMap, reverseLookUp


@numba.jit(nopython=True, parallel=True) #parallel speeds up computation only over very large matrices
# M is a mxn matrix binary matrix
# all elements in M should be uint8
def gf2elim(M):

    m,n = M.shape

    i=0
    j=0

    while i < m and j < n:
        # find value and index of largest element in remainder of column j
        k = np.argmax(M[i:, j])  +i

        # swap rows
        #M[[k, i]] = M[[i, k]] this doesn't work with numba
        temp = np.copy(M[k])
        M[k] = M[i]
        M[i] = temp


        aijn = M[i, j:]

        col = np.copy(M[:, j]) #make a copy otherwise M will be directly affected

        col[i] = 0 #avoid xoring pivot row with itself

        flip = np.outer(col, aijn)

        M[:, j:] = M[:, j:] ^ flip

        i += 1
        j +=1

    return M

if __name__ == '__main__':
    crack()